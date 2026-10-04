"""Превращает черновик LLM в ответ по контракту и страхует от ошибок модели."""

import logging
import re

from app.domain.assessment import (
    URGENCY_ORDER,
    AssessmentRequest,
    AssessmentResponse,
    LLMDraft,
    Recommendation,
)
from app.domain.catalog import CatalogItem
from app.domain.errors import InvalidDraftError
from app.domain.rules import Markers
from app.service.knowledge import KnowledgeBase

log = logging.getLogger(__name__)

# Пациенту нельзя показывать диагнозы, находки и пугающие формулировки.
_FORBIDDEN_FOR_PATIENT = re.compile(
    r"\bрак\w*|опухол|злокачеств|онко|метастаз|новообразован|карцином|малигн"
    r"|туберкул|BI-?RADS|пневмоторакс|гидроторакс|выпот|очаг|образовани\w*"
    r"|патолог|диагноз|подозрени",
    re.IGNORECASE,
)
_SAFE_PATIENT_TEXT = (
    "Врач рекомендует вам эту консультацию или обследование. Запишитесь, "
    "а специалист клиники подробно ответит на ваши вопросы."
)


def finalize(
    req: AssessmentRequest,
    markers: Markers,
    draft: LLMDraft,
    kb: KnowledgeBase,
) -> AssessmentResponse:
    already_recommended = {s.code for s in req.current_recommendations if s.code}
    booked = {a.service_code for a in req.appointments if a.service_code}

    recommendations: list[Recommendation] = []
    seen: set[str] = set()
    for rec in draft.recommendations:
        code = rec.service_code if rec.service_code in kb.catalog else ""
        if code and code in already_recommended:
            continue

        name = kb.catalog[code].name if code else rec.service_name.strip()
        if not name:
            log.warning("dropping recommendation without service name")
            continue

        key = code or name.lower()
        if key in seen:
            continue
        seen.add(key)

        recommendations.append(
            Recommendation(
                service_code=code,
                service_name=name,
                importance=rec.importance,
                rationale=rec.rationale.strip(),
                guideline_ref=kb.guideline_ref(rec.guideline_ref),
                patient_text=_safe_patient_text(rec.patient_text, kb.catalog.get(code)),
                already_booked=bool(code) and code in booked,
            )
        )

    # Сначала «крайне важные», внутри группы — порядок, который предложила LLM.
    recommendations.sort(key=lambda r: r.importance != "high")

    urgency = max(
        draft.urgency,
        markers.urgency_floor,
        key=URGENCY_ORDER.index,
    )
    if urgency != draft.urgency:
        log.info(
            "urgency raised by rules: %s -> %s (%s)",
            draft.urgency,
            urgency,
            ", ".join(markers.reasons),
        )
    # Пустой список допустим, если всё уже рекомендовано ранее. Но если модель
    # нашла повод для срочности и не предложила ни одного шага — это сбой,
    # main-B2B повторит запрос.
    if urgency != "normal" and not draft.recommendations:
        raise InvalidDraftError(f"urgency {urgency} without recommendations")

    return AssessmentResponse(
        urgency=urgency,
        catalog_version=kb.catalog_version,
        guidelines_version=kb.guidelines_version,
        recommendations=recommendations,
    )


def _safe_patient_text(text: str, item: CatalogItem | None) -> str:
    text = text.strip()
    if text and not _FORBIDDEN_FOR_PATIENT.search(text):
        return text
    # Замена: описание услуги из МИС (оно и так показывается пациенту), иначе шаблон.
    if item and item.description and not _FORBIDDEN_FOR_PATIENT.search(item.description):
        return item.description
    return _SAFE_PATIENT_TEXT
