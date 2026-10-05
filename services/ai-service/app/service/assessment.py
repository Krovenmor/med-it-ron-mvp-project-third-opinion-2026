"""Сценарий оценки заключения: база знаний → правила → LLM → постобработка."""

import logging

from app.domain import rules
from app.domain.assessment import AssessmentRequest, AssessmentResponse, Recommendation
from app.service.knowledge import KnowledgeProvider
from app.service.ports import LLM
from app.service.postprocess import finalize

log = logging.getLogger(__name__)


class AssessmentService:
    def __init__(self, knowledge: KnowledgeProvider, llm: LLM) -> None:
        self._knowledge = knowledge
        self._llm = llm

    async def assess(self, req: AssessmentRequest) -> AssessmentResponse:
        kb = await self._knowledge.get()
        markers = rules.extract(req)
        draft = await self._llm.assess(req, markers, kb)
        result = finalize(req, markers, draft, kb)
        log.info(
            "assessed case_id=%s urgency=%s catalog=%s recommendations=[%s]",
            req.case_id,
            result.urgency,
            kb.catalog_version,
            ", ".join(_summary(r) for r in result.recommendations),
        )
        return result


def _summary(rec: Recommendation) -> str:
    """Код (или название, если услуги нет в клинике), важность и отметка «уже записан» — без текстов."""
    label = rec.service_code or f"нет в клинике: {rec.service_name}"
    return f"{label} {rec.importance}{' booked' if rec.already_booked else ''}"
