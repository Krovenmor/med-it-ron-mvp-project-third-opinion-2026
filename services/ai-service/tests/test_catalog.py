import asyncio
import json

import httpx
import pytest

from app.domain.errors import CatalogUnavailableError
from app.infra.mis import MISCatalog, build_catalog
from conftest import MIS_SERVICES


class FakeMIS:
    """Фейковая МИС: отвечает на GET /api/v1/services и считает запросы."""

    def __init__(self, services=MIS_SERVICES, status: int = 200) -> None:
        self.services = services
        self.status = status
        self.requests: list[httpx.Request] = []

    def handler(self, request: httpx.Request) -> httpx.Response:
        self.requests.append(request)
        if self.status != 200:
            return httpx.Response(self.status, json={"error": "internal error"})
        return httpx.Response(200, json={"services": self.services})

    def catalog(self, ttl: float = 300) -> MISCatalog:
        client = httpx.AsyncClient(transport=httpx.MockTransport(self.handler))
        return MISCatalog("http://mis:8081/", timeout=5, ttl=ttl, client=client)


def test_fetches_services_from_mis():
    mis = FakeMIS()
    catalog = asyncio.run(mis.catalog().get())

    assert str(mis.requests[0].url) == "http://mis:8081/api/v1/services"
    assert set(catalog.items) == {s["code"] for s in MIS_SERVICES}
    item = catalog.items["ONC-MAMMO-CONSULT"]
    assert item.name == "Консультация онколога-маммолога"
    assert item.description.startswith("Маммолог")
    assert catalog.version.startswith("mis-")


def test_cached_within_ttl():
    mis = FakeMIS()
    source = mis.catalog(ttl=300)

    async def twice():
        await source.get()
        await source.get()

    asyncio.run(twice())
    assert len(mis.requests) == 1


def test_refreshed_after_ttl_and_version_follows_content():
    mis = FakeMIS()
    source = mis.catalog(ttl=0)

    async def scenario():
        first = await source.get()
        mis.services = MIS_SERVICES[:1]
        second = await source.get()
        return first, second

    first, second = asyncio.run(scenario())
    assert len(mis.requests) == 2
    assert first.version != second.version
    assert set(second.items) == {"ONC-MAMMO-CONSULT"}


def test_keeps_last_catalog_when_mis_fails():
    mis = FakeMIS()
    source = mis.catalog(ttl=0)

    async def scenario():
        first = await source.get()
        mis.status = 500
        return first, await source.get()

    first, second = asyncio.run(scenario())
    assert second.version == first.version
    assert second.items == first.items


@pytest.mark.parametrize(
    "mis",
    [FakeMIS(status=500), FakeMIS(services=[])],
    ids=["mis-error", "empty-catalog"],
)
def test_unavailable_without_any_catalog(mis):
    with pytest.raises(CatalogUnavailableError):
        asyncio.run(mis.catalog().get())


def test_connection_error_is_unavailable():
    def refuse(request):
        raise httpx.ConnectError("connection refused", request=request)

    client = httpx.AsyncClient(transport=httpx.MockTransport(refuse))
    with pytest.raises(CatalogUnavailableError):
        asyncio.run(MISCatalog("http://mis:8081", 5, 300, client=client).get())


def test_invalid_entries_are_skipped():
    catalog = build_catalog([*MIS_SERVICES[:1], {"code": "", "name": "Без кода"}, {"code": "NO-NAME"}])
    assert set(catalog.items) == {"ONC-MAMMO-CONSULT"}


def test_version_does_not_depend_on_order():
    assert build_catalog(MIS_SERVICES).version == build_catalog(list(reversed(MIS_SERVICES))).version


def test_matches_mis_demo_contract_example():
    # пример ответа из contracts/mis-demo.md
    body = json.loads(
        '{"services": [{"code": "MRI-BREAST", "name": "МРТ молочных желез с контрастированием",'
        ' "description": "МРТ даст врачу более подробное изображение для точного плана."}]}'
    )
    assert build_catalog(body["services"]).items["MRI-BREAST"].name.startswith("МРТ")
