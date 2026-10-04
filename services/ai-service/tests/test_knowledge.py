import asyncio

from app.infra.mis import build_catalog
from app.service.knowledge import KnowledgeProvider, build_knowledge
from conftest import MIS_SERVICES, FakeCatalog


def test_guidelines_loaded(guidelines):
    assert guidelines.version
    assert "MG-BIRADS-4-5" in guidelines.fragments
    assert guidelines.instructions


def test_output_schema_restricts_codes_to_mis_catalog(kb, mis_catalog):
    item = kb.output_schema["properties"]["recommendations"]["items"]["properties"]
    assert set(item["service_code"]["enum"]) == set(mis_catalog.items) | {""}
    assert set(item["guideline_ref"]["enum"]) == set(kb.fragments) | {""}


def test_reference_text_lists_mis_services(kb):
    assert "ONC-MAMMO-CONSULT | Консультация онколога-маммолога | Маммолог обсудит" in kb.reference_text
    assert "[MG-BIRADS-4-5]" in kb.reference_text


def test_reference_text_is_deterministic(guidelines):
    shuffled = build_catalog(list(reversed(MIS_SERVICES)))
    assert build_knowledge(guidelines, shuffled).reference_text == build_knowledge(
        guidelines, build_catalog(MIS_SERVICES)
    ).reference_text


def test_knowledge_rebuilt_only_when_catalog_changes(guidelines):
    catalog = FakeCatalog(build_catalog(MIS_SERVICES))
    provider = KnowledgeProvider(guidelines, catalog)

    first = asyncio.run(provider.get())
    assert asyncio.run(provider.get()) is first

    catalog.catalog = build_catalog(MIS_SERVICES[:2])
    second = asyncio.run(provider.get())
    assert second is not first
    assert set(second.catalog) == {"ONC-MAMMO-CONSULT", "US-BREAST"}
