"""Справочник услуг клиники. Источник — МИС (infra/mis.py), своей копии у сервиса нет."""

from dataclasses import dataclass


@dataclass(frozen=True)
class CatalogItem:
    code: str
    name: str
    description: str  # понятное пациенту описание услуги


@dataclass(frozen=True)
class Catalog:
    # В МИС нет версии справочника — версией служит хеш его содержимого.
    version: str
    items: dict[str, CatalogItem]
    fetched_at: float = 0.0
