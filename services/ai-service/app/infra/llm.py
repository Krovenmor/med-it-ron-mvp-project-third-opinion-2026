"""Общее для адаптеров LLM: формат запроса, разбор ответа, ограничение параллельности."""

import asyncio
import json
import logging
import re

from pydantic import ValidationError

from app.domain.assessment import AssessmentRequest, LLMDraft
from app.domain.errors import LLMError
from app.domain.rules import Markers
from app.service.knowledge import KnowledgeBase
from app.service.ports import LLM

log = logging.getLogger(__name__)

_CODE_FENCE = re.compile(r"^```(?:json)?\s*(.*?)\s*```$", re.DOTALL)
_LOGGED_DRAFT_CHARS = 2000


class InvalidOutputError(LLMError):
    """Модель ответила, но не JSON по схеме. Такой сбой случаен — запрос имеет смысл повторить."""


class ConcurrencyLimited:
    """Ограничивает число одновременных запросов к LLM, остальные ждут очереди.

    main-B2B шлёт оценки параллельно (WORKER_COUNT), а тариф LLM может разрешать
    меньше одновременных запросов: без ограничения часть запросов получила бы 429.
    """

    def __init__(self, inner: LLM, limit: int) -> None:
        self._inner = inner
        self._semaphore = asyncio.Semaphore(limit)

    async def assess(self, req: AssessmentRequest, markers: Markers, kb: KnowledgeBase) -> LLMDraft:
        async with self._semaphore:
            return await self._inner.assess(req, markers, kb)

    async def aclose(self) -> None:
        await self._inner.aclose()


def json_format_instructions(schema: dict) -> str:
    """Описание формата ответа для режима, когда схему модели передают текстом в промпте."""
    return (
        "# Формат ответа\n\n"
        "Ответь одним JSON-объектом строго по этой JSON-схеме: без пояснений до и после, "
        "без markdown, только допустимые значения из enum.\n\n"
        + json.dumps(schema, ensure_ascii=False, sort_keys=True)
    )


def user_message(req: AssessmentRequest, markers: Markers) -> str:
    # case_id модели не нужен: он только для корреляции и логов.
    case = req.model_dump(mode="json", exclude={"case_id"})
    case["extracted_markers"] = markers.as_hint()
    return json.dumps(case, ensure_ascii=False, indent=2)


def parse_draft(text: str | None) -> LLMDraft:
    if not text:
        raise InvalidOutputError("llm returned empty content")
    text = text.strip()
    # Если режим JSON-схемы не сработал, модель может обернуть ответ в ```json ... ```.
    fenced = _CODE_FENCE.match(text)
    if fenced:
        text = fenced.group(1)
    try:
        return LLMDraft.model_validate_json(text)
    except ValidationError as e:
        # Ответ модели обезличен (без ФИО и ID), поэтому для разбора ошибок его можно писать в лог.
        log.warning("llm returned invalid draft: %s\n%s", _describe(e), text[:_LOGGED_DRAFT_CHARS])
        raise InvalidOutputError(f"llm returned invalid draft: {_describe(e)}") from e


def _describe(e: ValidationError, limit: int = 3) -> str:
    """Какие поля модель заполнила не по схеме: `recommendations.0.importance: ... (got 'medium')`."""
    parts = []
    for err in e.errors()[:limit]:
        loc = ".".join(str(x) for x in err["loc"]) or "<root>"
        got = repr(err.get("input"))
        parts.append(f"{loc}: {err['msg']} (got {got[:60]})")
    if e.error_count() > limit:
        parts.append(f"... and {e.error_count() - limit} more")
    return "; ".join(parts)
