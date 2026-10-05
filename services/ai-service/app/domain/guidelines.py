"""Выдержки из клинических рекомендаций МЗ РФ. Файлы — knowledge/guidelines (infra/guidelines.py)."""

from dataclasses import dataclass


@dataclass(frozen=True)
class GuidelineFragment:
    id: str
    ref: str  # ссылка, которую увидит врач
    text: str


@dataclass(frozen=True)
class Guidelines:
    version: str
    fragments: dict[str, GuidelineFragment]
    instructions: str  # системная инструкция для LLM
