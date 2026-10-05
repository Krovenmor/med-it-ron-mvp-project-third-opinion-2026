"""Сборка сервиса из зависимостей — аналог internal/di в Go-сервисах проекта.

Запуск: python -m app (адрес — HTTP_ADDR). Тесты передают в create_app фейки вместо МИС и LLM.
"""

import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI

from app.config import Settings, load_settings
from app.domain.errors import CatalogUnavailableError
from app.domain.guidelines import Guidelines
from app.infra.gigachat import GigaChatAssessor
from app.infra.guidelines import load_guidelines
from app.infra.llm import ConcurrencyLimited
from app.infra.mis import MISCatalog
from app.service.assessment import AssessmentService
from app.service.knowledge import KnowledgeProvider
from app.service.ports import LLM, CatalogSource
from app.transport.http import new_app

log = logging.getLogger("ai-service")


def create_app(
    settings: Settings | None = None,
    guidelines: Guidelines | None = None,
    catalog: CatalogSource | None = None,
    llm: LLM | None = None,
) -> FastAPI:
    settings = settings or load_settings()
    logging.basicConfig(
        level=settings.log_level,
        format="%(asctime)s %(levelname)s %(name)s: %(message)s",
    )
    guidelines = guidelines or load_guidelines(settings.knowledge_dir)
    catalog = catalog or MISCatalog(settings.mis_url, settings.mis_timeout, settings.catalog_ttl)
    llm = llm or create_llm(settings)
    knowledge = KnowledgeProvider(guidelines, catalog)
    service = AssessmentService(knowledge, llm)
    log.info(
        "starting: llm=%s, mis=%s, guidelines=%s (%d fragments)",
        settings.llm_model,
        settings.mis_url,
        guidelines.version,
        len(guidelines.fragments),
    )

    @asynccontextmanager
    async def lifespan(_: FastAPI):
        # Справочник загружаем заранее, но запуск от МИС не зависит:
        # если она ещё не поднялась, справочник догрузится при первой оценке.
        try:
            await knowledge.get()
        except CatalogUnavailableError as e:
            log.warning("mis catalog not loaded yet: %s", e)
        yield
        await llm.aclose()
        await catalog.aclose()

    def health() -> dict:
        loaded = catalog.cached
        return {
            "status": "ok",
            "llm": settings.llm_model,
            "catalog": {"version": loaded.version, "services": len(loaded.items)} if loaded else None,
            "guidelines_version": guidelines.version,
        }

    return new_app(service, health, lifespan)


def create_llm(settings: Settings) -> LLM:
    llm: LLM = GigaChatAssessor(settings)
    if settings.llm_max_concurrency > 0:
        llm = ConcurrencyLimited(llm, settings.llm_max_concurrency)
    return llm
