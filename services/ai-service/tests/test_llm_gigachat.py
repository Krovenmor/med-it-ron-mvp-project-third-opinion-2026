import asyncio
import json

import pytest
from gigachat.exceptions import ResponseError
from gigachat.models import ChatCompletion

from app.config import load_settings
from app.domain import rules
from app.domain.errors import LLMError
from app.infra.gigachat import GigaChatAssessor
from app.infra.llm import ConcurrencyLimited
from app.main import create_llm
from conftest import make_request, parse_request

DRAFT = {
    "urgency": "priority",
    "recommendations": [
        {
            "service_code": "ONC-MAMMO-CONSULT",
            "service_name": "Консультация онколога-маммолога",
            "importance": "high",
            "rationale": "BI-RADS 4.",
            "guideline_ref": "MG-BIRADS-4-5",
            "patient_text": "Рекомендуем консультацию маммолога.",
        }
    ],
}


def completion(content: str, finish_reason: str = "stop") -> ChatCompletion:
    return ChatCompletion.model_validate(
        {
            "choices": [
                {
                    "message": {"role": "assistant", "content": content},
                    "index": 0,
                    "finish_reason": finish_reason,
                }
            ],
            "created": 1,
            "model": "GigaChat-2-Max",
            "usage": {"prompt_tokens": 3000, "completion_tokens": 300, "total_tokens": 3300},
            "object": "chat.completion",
        }
    )


class FakeGigaChat:
    def __init__(self, result) -> None:
        self.result = result
        self.chats = []

    async def achat(self, chat):
        self.chats.append(chat)
        if isinstance(self.result, Exception):
            raise self.result
        return self.result

    async def aclose(self) -> None:
        pass


def run(assessor: GigaChatAssessor, kb):
    req = parse_request(make_request("Заключение: BI-RADS 4.", "MG"))
    return asyncio.run(assessor.assess(req, rules.extract(req), kb))


@pytest.fixture
def settings():
    return load_settings()


def test_request_shape_and_parsing(settings, kb):
    fake = FakeGigaChat(completion(json.dumps(DRAFT, ensure_ascii=False)))
    draft = run(GigaChatAssessor(settings, client=fake), kb)

    assert draft.urgency == "priority"
    assert draft.recommendations[0].service_code == "ONC-MAMMO-CONSULT"

    chat = fake.chats[0]
    system, user = chat.messages
    assert system.role == "system" and kb.reference_text in system.content
    assert user.role == "user" and "case_id" not in user.content
    assert chat.temperature == settings.llm_temperature
    # режим prompt (по умолчанию): схема с enum по кодам справочника — в системном промпте
    assert settings.llm_output_mode == "prompt"
    assert chat.response_format is None
    assert "# Формат ответа" in system.content
    assert '"ONC-MAMMO-CONSULT"' in system.content.split("# Формат ответа")[1]


def test_json_schema_mode_uses_response_format(settings, kb):
    from dataclasses import replace

    fake = FakeGigaChat(completion(json.dumps(DRAFT, ensure_ascii=False)))
    run(GigaChatAssessor(replace(settings, llm_output_mode="json_schema"), client=fake), kb)
    chat = fake.chats[0]
    assert chat.response_format.schema_ == kb.output_schema
    assert chat.response_format.strict is True
    assert "# Формат ответа" not in chat.messages[0].content


def test_code_fenced_json_is_accepted(settings, kb):
    content = "```json\n" + json.dumps(DRAFT, ensure_ascii=False) + "\n```"
    draft = run(GigaChatAssessor(settings, client=FakeGigaChat(completion(content))), kb)
    assert draft.urgency == "priority"


@pytest.mark.parametrize(
    ("result", "message"),
    [
        (completion("{}", finish_reason="length"), "finished with length"),
        (completion("{}", finish_reason="blacklist"), "finished with blacklist"),
        (completion("не json"), "invalid draft"),
        (completion(""), "empty content"),
        (ResponseError("https://api", 429, b"too many requests", None), "responded 429"),
    ],
)
def test_failures_become_llm_error(settings, kb, result, message):
    with pytest.raises(LLMError, match=message):
        run(GigaChatAssessor(settings, client=FakeGigaChat(result)), kb)


def test_factory_builds_gigachat_with_concurrency_limit(monkeypatch):
    monkeypatch.setenv("LLM_MAX_CONCURRENCY", "1")
    llm = create_llm(load_settings())
    assert isinstance(llm, ConcurrencyLimited)
    assert isinstance(llm._inner, GigaChatAssessor)




def test_invalid_draft_names_the_field(settings, kb):
    bad = json.loads(json.dumps(DRAFT))
    bad["recommendations"][0]["importance"] = "medium"
    content = json.dumps(bad, ensure_ascii=False)
    with pytest.raises(LLMError, match=r"recommendations\.0\.importance: .*got 'medium'"):
        run(GigaChatAssessor(settings, client=FakeGigaChat(completion(content))), kb)


def test_retries_once_on_invalid_json(settings, kb):
    good = completion(json.dumps(DRAFT, ensure_ascii=False))
    broken = completion('{"urgency": "priority", "recommendations": [{ Отец "service_code": 1}]}')

    class Flaky(FakeGigaChat):
        def __init__(self):
            super().__init__(None)
            self.results = [broken, good]

        async def achat(self, chat):
            self.chats.append(chat)
            return self.results.pop(0)

    fake = Flaky()
    draft = run(GigaChatAssessor(settings, client=fake), kb)
    assert draft.urgency == "priority"
    assert len(fake.chats) == 2


def test_gives_up_after_second_invalid_json(settings, kb):
    fake = FakeGigaChat(completion("не json"))
    with pytest.raises(LLMError, match="invalid draft"):
        run(GigaChatAssessor(settings, client=fake), kb)
    assert len(fake.chats) == 2


def test_no_retry_when_answer_truncated(settings, kb):
    fake = FakeGigaChat(completion("{}", finish_reason="length"))
    with pytest.raises(LLMError, match="length"):
        run(GigaChatAssessor(settings, client=fake), kb)
    assert len(fake.chats) == 1


def test_repetition_penalty_disabled(settings, kb):
    fake = FakeGigaChat(completion(json.dumps(DRAFT, ensure_ascii=False)))
    run(GigaChatAssessor(settings, client=fake), kb)
    assert fake.chats[0].repetition_penalty == 1.0
