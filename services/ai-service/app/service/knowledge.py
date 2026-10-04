"""База знаний для LLM: справочник услуг из МИС + выдержки из клинических рекомендаций.

Из них собираются текст для системного промпта и JSON-схема ответа модели. Всё собирается
детерминированно (стабильный порядок, без дат), чтобы у провайдера LLM срабатывал кэш
одинакового начала промпта.
"""

from __future__ import annotations

from dataclasses import dataclass
from typing import TYPE_CHECKING

from app.domain.catalog import Catalog, CatalogItem
from app.domain.guidelines import GuidelineFragment, Guidelines

if TYPE_CHECKING:
    from app.service.ports import CatalogSource


@dataclass(frozen=True)
class KnowledgeBase:
    catalog_version: str
    guidelines_version: str
    catalog: dict[str, CatalogItem]
    fragments: dict[str, GuidelineFragment]
    instructions: str
    reference_text: str
    output_schema: dict

    def guideline_ref(self, fragment_id: str) -> str:
        fragment = self.fragments.get(fragment_id)
        return fragment.ref if fragment else ""


class KnowledgeProvider:
    """Отдаёт актуальную базу знаний; пересобирает её, только когда меняется справочник МИС."""

    def __init__(self, guidelines: Guidelines, catalog: CatalogSource) -> None:
        self._guidelines = guidelines
        self._catalog = catalog
        self._kb: KnowledgeBase | None = None

    async def get(self) -> KnowledgeBase:
        catalog = await self._catalog.get()
        if self._kb is None or self._kb.catalog_version != catalog.version:
            self._kb = build_knowledge(self._guidelines, catalog)
        return self._kb


def build_knowledge(guidelines: Guidelines, catalog: Catalog) -> KnowledgeBase:
    return KnowledgeBase(
        catalog_version=catalog.version,
        guidelines_version=guidelines.version,
        catalog=catalog.items,
        fragments=guidelines.fragments,
        instructions=guidelines.instructions,
        reference_text=_render_reference(catalog.items, guidelines.fragments),
        output_schema=_output_schema(sorted(catalog.items), sorted(guidelines.fragments)),
    )


def _render_reference(catalog: dict[str, CatalogItem], fragments: dict[str, GuidelineFragment]) -> str:
    lines = ["# Каталог услуг клиники", "", "Код | Название | Описание для пациента", "--- | --- | ---"]
    lines += [f"{c.code} | {c.name} | {c.description or '—'}" for c in catalog.values()]
    lines += ["", "# Выдержки из клинических рекомендаций", ""]
    for f in fragments.values():
        lines += [f"## [{f.id}] {f.ref}", f.text, ""]
    return "\n".join(lines).strip()


def _output_schema(service_codes: list[str], fragment_ids: list[str]) -> dict:
    """JSON-схема ответа LLM. enum не даёт модели выдумать код услуги или ссылку на КР."""
    return {
        "type": "object",
        "properties": {
            "urgency": {"type": "string", "enum": ["normal", "planned", "priority", "emergency"]},
            "recommendations": {
                "type": "array",
                "items": {
                    "type": "object",
                    "properties": {
                        "service_code": {"type": "string", "enum": [*service_codes, ""]},
                        "service_name": {"type": "string"},
                        "importance": {"type": "string", "enum": ["high", "low"]},
                        "rationale": {
                            "type": "string",
                            "description": "Для врача: какая находка, зачем услуга, на чём основано",
                        },
                        "guideline_ref": {"type": "string", "enum": [*fragment_ids, ""]},
                        "patient_text": {
                            "type": "string",
                            "description": (
                                "Для пациента, 1–2 предложения: зачем нужен шаг. Без диагнозов и находок: "
                                "нельзя «образование», «очаг», «рак», «опухоль», BI-RADS, размеры"
                            ),
                        },
                    },
                    "required": [
                        "service_code",
                        "service_name",
                        "importance",
                        "rationale",
                        "guideline_ref",
                        "patient_text",
                    ],
                    "additionalProperties": False,
                },
            },
        },
        "required": ["urgency", "recommendations"],
        "additionalProperties": False,
    }
