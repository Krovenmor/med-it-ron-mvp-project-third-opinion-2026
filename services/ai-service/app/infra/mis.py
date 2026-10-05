"""Справочник услуг клиники из МИС: GET {MIS_URL}/api/v1/services (contracts/mis-demo.md).

Рекомендовать можно только услуги, которые есть в МИС: по их кодам b2c-service ищет
слоты и записи, а b2b-service сверяет визиты.

Справочник кэшируется на CATALOG_TTL секунд. Если обновить его не удалось, используется
последняя полученная версия; если МИС недоступна с самого старта — CatalogUnavailableError,
оценка отвечает 502, и main-B2B повторит запрос позже.
"""

import asyncio
import hashlib
import json
import logging
import time
from dataclasses import replace

import httpx

from app.domain.catalog import Catalog, CatalogItem
from app.domain.errors import CatalogUnavailableError

log = logging.getLogger(__name__)


def build_catalog(services: list[dict]) -> Catalog:
    """Собирает справочник из ответа МИС; записи без кода или названия пропускаются."""
    items: dict[str, CatalogItem] = {}
    for s in sorted(services, key=lambda s: str(s.get("code", ""))):
        code, name = str(s.get("code", "")).strip(), str(s.get("name", "")).strip()
        if not code or not name:
            log.warning("skipping mis service without code or name: %r", s)
            continue
        items[code] = CatalogItem(code=code, name=name, description=str(s.get("description", "")).strip())
    if not items:
        raise ValueError("mis returned an empty service catalog")

    digest = hashlib.sha256(
        json.dumps([[i.code, i.name, i.description] for i in items.values()], ensure_ascii=False).encode()
    ).hexdigest()[:10]
    return Catalog(version=f"mis-{digest}", items=items, fetched_at=time.monotonic())


class MISCatalog:
    def __init__(
        self,
        mis_url: str,
        timeout: float,
        ttl: float,
        client: httpx.AsyncClient | None = None,
    ) -> None:
        self._url = mis_url.rstrip("/") + "/api/v1/services"
        self._ttl = ttl
        self._client = client or httpx.AsyncClient(timeout=timeout)
        self._cached: Catalog | None = None
        self._lock = asyncio.Lock()

    @property
    def cached(self) -> Catalog | None:
        return self._cached

    async def get(self) -> Catalog:
        if self._is_fresh():
            return self._cached  # type: ignore[return-value]
        async with self._lock:
            if self._is_fresh():
                return self._cached  # type: ignore[return-value]
            try:
                catalog = await self._fetch()
            except (httpx.HTTPError, ValueError) as e:
                if self._cached is None:
                    raise CatalogUnavailableError(f"{self._url}: {e}") from e
                log.warning("mis catalog refresh failed, keeping %s: %s", self._cached.version, e)
                # Следующая попытка — через TTL, а не на каждом запросе.
                self._cached = replace(self._cached, fetched_at=time.monotonic())
                return self._cached

            if self._cached is None or self._cached.version != catalog.version:
                log.info("mis catalog loaded: %s, %d services", catalog.version, len(catalog.items))
            self._cached = catalog
            return catalog

    async def aclose(self) -> None:
        await self._client.aclose()

    def _is_fresh(self) -> bool:
        return self._cached is not None and time.monotonic() - self._cached.fetched_at < self._ttl

    async def _fetch(self) -> Catalog:
        resp = await self._client.get(self._url)
        resp.raise_for_status()
        services = resp.json().get("services")
        if not isinstance(services, list):
            raise ValueError("response has no 'services' array")
        return build_catalog(services)
