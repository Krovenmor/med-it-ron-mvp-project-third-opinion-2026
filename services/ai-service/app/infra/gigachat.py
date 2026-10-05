"""Адаптер LLM на GigaChat (Сбер), реализует порт service.ports.LLM.

Доступ SDK берёт из переменных окружения GIGACHAT_*:
- GIGACHAT_CREDENTIALS     — Authorization key из личного кабинета developers.sber.ru;
- GIGACHAT_SCOPE           — GIGACHAT_API_PERS (физлица), GIGACHAT_API_B2B / GIGACHAT_API_CORP (юрлица);
- GIGACHAT_CA_BUNDLE_FILE  — сертификат НУЦ Минцифры (иначе ошибка SSL); путь разбирает app/config.py;
- GIGACHAT_VERIFY_SSL_CERTS=False — только для локальной отладки.
OAuth-токен (живёт 30 минут) SDK получает и обновляет сам.
"""

import logging
import time
from typing import Any

import httpx
from gigachat import GigaChat
from gigachat.exceptions import GigaChatException, ResponseError
from gigachat.models import Chat, ChatCompletion, Messages, MessagesRole
from gigachat.models.response_format import JsonSchemaResponseFormat

from app.config import Settings
from app.domain.assessment import AssessmentRequest, LLMDraft
from app.domain.errors import LLMError
from app.domain.rules import Markers
from app.infra.llm import InvalidOutputError, json_format_instructions, parse_draft, user_message
from app.service.knowledge import KnowledgeBase

log = logging.getLogger(__name__)

# Невалидный JSON у GigaChat случаен, поэтому одна повторная попытка обычно спасает ответ.
# Повторяем, только если первая попытка уложилась в половину LLM_TIMEOUT: иначе не успеем
# до таймаута main-B2B (90 с) — пусть лучше он сам повторит позже.
_ATTEMPTS = 2


class GigaChatAssessor:
    def __init__(self, settings: Settings, client: Any = None) -> None:
        self._settings = settings
        # Сетевые повторы делает main-B2B, здесь их не дублируем.
        ca_bundle = settings.gigachat_ca_bundle_file
        self._client = client or GigaChat(
            model=settings.llm_model,
            timeout=settings.llm_timeout,
            max_retries=0,
            # None — SDK возьмёт GIGACHAT_CA_BUNDLE_FILE из окружения (там пусто) и системные сертификаты.
            ca_bundle_file=str(ca_bundle) if ca_bundle else None,
        )

    async def assess(self, req: AssessmentRequest, markers: Markers, kb: KnowledgeBase) -> LLMDraft:
        system = f"{kb.instructions}\n\n{kb.reference_text}"
        response_format = None
        if self._settings.llm_output_mode == "json_schema":
            # beta-режим GigaChat: генерацию не ограничивает и на практике вставляет мусор
            # в отступы перед ключами JSON. Оставлен, чтобы сравнить, когда Сбер его доработает.
            response_format = JsonSchemaResponseFormat(schema=kb.output_schema, strict=True)
        else:
            system = f"{system}\n\n{json_format_instructions(kb.output_schema)}"

        chat = Chat(
            messages=[
                # У GigaChat одно системное сообщение, и оно идёт первым. Пока справочник
                # МИС не меняется, оно одинаково во всех запросах — часть токенов
                # GigaChat может отдать из кэша.
                Messages(role=MessagesRole.SYSTEM, content=system),
                Messages(role=MessagesRole.USER, content=user_message(req, markers)),
            ],
            temperature=self._settings.llm_temperature,
            repetition_penalty=self._settings.llm_repetition_penalty,
            max_tokens=self._settings.llm_max_tokens,
            response_format=response_format,
        )

        started = time.monotonic()
        for attempt in range(1, _ATTEMPTS + 1):
            completion = await self._complete(chat, req.case_id, attempt)
            try:
                return self._draft(completion)
            except InvalidOutputError:
                elapsed = time.monotonic() - started
                if attempt == _ATTEMPTS or elapsed > self._settings.llm_timeout / 2:
                    raise
                log.warning("case_id=%s: invalid draft, retrying (%.1fs spent)", req.case_id, elapsed)
        raise AssertionError("unreachable")

    async def aclose(self) -> None:
        await self._client.aclose()

    async def _complete(self, chat: Chat, case_id: str, attempt: int) -> ChatCompletion:
        started = time.monotonic()
        try:
            completion = await self._client.achat(chat)
        except ResponseError as e:
            log.error("gigachat responded %s: %r", e.status_code, e.content)
            raise LLMError(f"gigachat responded {e.status_code}") from e
        except (GigaChatException, httpx.HTTPError) as e:
            log.error("gigachat request failed: %s", e)
            raise LLMError(f"gigachat request failed: {type(e).__name__}") from e

        choice = completion.choices[0] if completion.choices else None
        usage = completion.usage
        log.info(
            "llm done case_id=%s attempt=%d model=%s finish=%s %.1fs in=%s cached=%s out=%s",
            case_id,
            attempt,
            completion.model,
            choice.finish_reason if choice else None,
            time.monotonic() - started,
            usage.prompt_tokens,
            usage.precached_prompt_tokens,
            usage.completion_tokens,
        )
        return completion

    @staticmethod
    def _draft(completion: ChatCompletion) -> LLMDraft:
        choice = completion.choices[0] if completion.choices else None
        if choice is None:
            raise LLMError("gigachat returned no choices")
        # stop — штатное завершение; length — упёрлись в max_tokens; blacklist — фильтр контента.
        # Их повтор не исправит, поэтому это LLMError, а не InvalidOutputError.
        if choice.finish_reason != "stop":
            raise LLMError(f"gigachat finished with {choice.finish_reason}")
        return parse_draft(choice.message.content)
