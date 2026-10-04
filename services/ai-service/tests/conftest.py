import pytest
from fastapi.testclient import TestClient

from app import config
from app.config import load_settings
from app.domain.assessment import AssessmentRequest, LLMDraft
from app.domain.catalog import Catalog
from app.domain.errors import CatalogUnavailableError
from app.infra.guidelines import load_guidelines
from app.infra.mis import build_catalog
from app.main import create_app
from app.service.knowledge import build_knowledge

# Тесты не читают локальный .env с ключами — только несекретный config/default.env.
config.ENV_FILES = (config.BASE_DIR / "config" / "default.env",)

# Фрагмент справочника mis-demo (ответ GET /api/v1/services, contracts/mis-demo.md).
MIS_SERVICES = [
    {
        "code": "ONC-MAMMO-CONSULT",
        "name": "Консультация онколога-маммолога",
        "description": "Маммолог обсудит с вами результаты обследования и дальнейшие шаги.",
    },
    {
        "code": "US-BREAST",
        "name": "УЗИ молочных желез",
        "description": "Ультразвуковое исследование поможет врачу получить больше информации.",
    },
    {
        "code": "CARDIO-CONSULT",
        "name": "Консультация кардиолога",
        "description": "Кардиолог проверит работу сердца и при необходимости скорректирует план.",
    },
    {
        "code": "LAB-BIOCHEM",
        "name": "Анализ крови (биохимия)",
        "description": "Биохимический анализ крови даст врачу дополнительную информацию.",
    },
    {
        "code": "PULM-CONSULT",
        "name": "Консультация пульмонолога",
        "description": "Пульмонолог оценит результаты исследования лёгких и подскажет дальнейшие шаги.",
    },
    {
        "code": "SURG-THOR-CONSULT",
        "name": "Консультация торакального хирурга",
        "description": "Хирург оценит результаты исследования и подскажет, нужна ли дополнительная помощь.",
    },
    {
        "code": "SCREENING-ANNUAL",
        "name": "Плановый профилактический осмотр",
        "description": "Ежегодный осмотр, чтобы вовремя заметить изменения.",
    },
]


def make_request(conclusion: str, modality: str = "DX", **overrides) -> dict:
    body = {
        "case_id": "39e61250-869e-4a3c-95a4-b16428e14332",
        "study": {
            "modality": modality,
            "body_site": "chest",
            "performed_at": "2026-10-01T09:00:00Z",
        },
        "conclusion": conclusion,
        "patient": {"age": 52, "sex": "female"},
        "current_recommendations": [],
        "visits": [],
        "appointments": [],
    }
    body.update(overrides)
    return body


def parse_request(body: dict) -> AssessmentRequest:
    return AssessmentRequest.model_validate(body)


class FakeAssessor:
    """Подменяет LLM: возвращает заранее заданный черновик или бросает ошибку."""

    def __init__(self) -> None:
        self.draft: LLMDraft | None = None
        self.error: Exception | None = None
        self.calls: list[tuple] = []

    async def assess(self, req, markers, kb) -> LLMDraft:
        self.calls.append((req, markers, kb))
        if self.error:
            raise self.error
        assert self.draft is not None, "test must set FakeAssessor.draft"
        return self.draft

    async def aclose(self) -> None:
        pass


class FakeCatalog:
    """Подменяет справочник МИС."""

    def __init__(self, catalog: Catalog | None) -> None:
        self.catalog = catalog

    @property
    def cached(self) -> Catalog | None:
        return self.catalog

    async def get(self) -> Catalog:
        if self.catalog is None:
            raise CatalogUnavailableError("mis is down")
        return self.catalog

    async def aclose(self) -> None:
        pass


@pytest.fixture(scope="session")
def guidelines():
    return load_guidelines(load_settings().knowledge_dir)


@pytest.fixture(scope="session")
def mis_catalog() -> Catalog:
    return build_catalog(MIS_SERVICES)


@pytest.fixture(scope="session")
def kb(guidelines, mis_catalog):
    return build_knowledge(guidelines, mis_catalog)


@pytest.fixture
def fake_llm() -> FakeAssessor:
    return FakeAssessor()


@pytest.fixture
def fake_catalog(mis_catalog) -> FakeCatalog:
    return FakeCatalog(mis_catalog)


@pytest.fixture
def client(guidelines, fake_catalog, fake_llm) -> TestClient:
    return TestClient(create_app(guidelines=guidelines, catalog=fake_catalog, llm=fake_llm))
