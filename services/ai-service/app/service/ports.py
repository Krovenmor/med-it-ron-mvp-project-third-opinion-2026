"""Порты сценария оценки: что сервису нужно от внешнего мира. Реализации — в infra."""

from typing import Protocol

from app.domain.assessment import AssessmentRequest, LLMDraft
from app.domain.catalog import Catalog
from app.domain.rules import Markers
from app.service.knowledge import KnowledgeBase


class LLM(Protocol):
    """Модель, которая по заключению и базе знаний предлагает черновик рекомендаций."""

    async def assess(self, req: AssessmentRequest, markers: Markers, kb: KnowledgeBase) -> LLMDraft: ...

    async def aclose(self) -> None: ...


class CatalogSource(Protocol):
    """Справочник услуг клиники."""

    @property
    def cached(self) -> Catalog | None: ...

    async def get(self) -> Catalog: ...

    async def aclose(self) -> None: ...
