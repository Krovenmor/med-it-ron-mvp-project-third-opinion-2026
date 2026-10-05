from conftest import make_request
from app.domain.assessment import DraftRecommendation, LLMDraft
from app.domain.errors import LLMError
from app.infra.llm import user_message


def rec(code: str, **kw) -> DraftRecommendation:
    defaults = {
        "service_code": code,
        "service_name": "",
        "importance": "high",
        "rationale": "обоснование",
        "guideline_ref": "",
        "patient_text": "Рекомендуем консультацию, чтобы обсудить результаты.",
    }
    return DraftRecommendation(**(defaults | kw))


def test_contract_response(client, fake_llm, kb):
    fake_llm.draft = LLMDraft(
        urgency="priority",
        recommendations=[
            rec("US-BREAST", importance="low", guideline_ref="MG-BIRADS-4-5"),
            rec("ONC-MAMMO-CONSULT", guideline_ref="MG-BIRADS-4-5"),
        ],
    )
    resp = client.post("/v1/assessments", json=make_request("BI-RADS 4", "MG"))

    assert resp.status_code == 200
    body = resp.json()
    assert body["urgency"] == "priority"
    assert body["catalog_version"] == kb.catalog_version
    assert body["guidelines_version"] == kb.guidelines_version
    # high идёт первым, название и ссылка подставлены из базы знаний
    first, second = body["recommendations"]
    assert first["service_code"] == "ONC-MAMMO-CONSULT"
    assert first["service_name"] == "Консультация онколога-маммолога"
    assert first["guideline_ref"] == kb.fragments["MG-BIRADS-4-5"].ref
    assert second["service_code"] == "US-BREAST"
    assert set(first) == {
        "service_code",
        "service_name",
        "importance",
        "rationale",
        "guideline_ref",
        "patient_text",
        "already_booked",
    }


def test_case_id_not_sent_to_llm(client, fake_llm):
    fake_llm.draft = LLMDraft(urgency="normal", recommendations=[rec("SCREENING-ANNUAL")])
    client.post("/v1/assessments", json=make_request("Без патологии"))
    req, markers, _ = fake_llm.calls[0]
    assert "case_id" not in user_message(req, markers)


def test_rules_raise_urgency_floor(client, fake_llm):
    fake_llm.draft = LLMDraft(urgency="planned", recommendations=[rec("SURG-THOR-CONSULT")])
    resp = client.post("/v1/assessments", json=make_request("Справа пневмоторакс."))
    assert resp.json()["urgency"] == "emergency"


def test_already_booked_and_current_recommendations(client, fake_llm):
    fake_llm.draft = LLMDraft(
        urgency="planned",
        recommendations=[rec("CARDIO-CONSULT"), rec("LAB-BIOCHEM"), rec("PULM-CONSULT")],
    )
    body = make_request(
        "Кардиомегалия.",
        current_recommendations=[{"code": "PULM-CONSULT", "name": "Консультация пульмонолога"}],
        appointments=[
            {
                "service_code": "LAB-BIOCHEM",
                "service_name": "Анализ крови (биохимия)",
                "scheduled_at": "2026-10-10T10:00:00Z",
            }
        ],
    )
    recs = {r["service_code"]: r for r in client.post("/v1/assessments", json=body).json()["recommendations"]}
    assert set(recs) == {"CARDIO-CONSULT", "LAB-BIOCHEM"}
    assert recs["LAB-BIOCHEM"]["already_booked"] is True
    assert recs["CARDIO-CONSULT"]["already_booked"] is False


def test_service_missing_in_clinic_and_duplicates(client, fake_llm):
    fake_llm.draft = LLMDraft(
        urgency="planned",
        recommendations=[
            rec("", service_name="Консультация эндокринолога"),
            rec("", service_name="консультация эндокринолога"),
            rec("CARDIO-CONSULT"),
            rec("CARDIO-CONSULT"),
            rec("", service_name="  "),
        ],
    )
    recs = client.post("/v1/assessments", json=make_request("Кардиомегалия.")).json()["recommendations"]
    assert [(r["service_code"], r["service_name"]) for r in recs] == [
        ("", "Консультация эндокринолога"),
        ("CARDIO-CONSULT", "Консультация кардиолога"),
    ]


def test_patient_text_without_diagnosis(client, fake_llm):
    fake_llm.draft = LLMDraft(
        urgency="priority",
        recommendations=[rec("ONC-MAMMO-CONSULT", patient_text="У вас подозрение на рак, BI-RADS 4.")],
    )
    text = client.post("/v1/assessments", json=make_request("BI-RADS 4", "MG")).json()[
        "recommendations"
    ][0]["patient_text"]
    assert "рак" not in text.lower()
    assert "BI-RADS" not in text
    # подставлено понятное пациенту описание услуги из справочника МИС
    assert text == "Маммолог обсудит с вами результаты обследования и дальнейшие шаги."


def test_unknown_fields_are_ignored(client, fake_llm):
    fake_llm.draft = LLMDraft(urgency="normal", recommendations=[rec("SCREENING-ANNUAL")])
    body = make_request("Без патологии", ai_result={"pathologyFlag": False})
    body["study"]["series_count"] = 3
    assert client.post("/v1/assessments", json=body).status_code == 200


def test_invalid_request_returns_error_json(client):
    body = make_request("текст")
    body["study"]["modality"] = "MR"
    resp = client.post("/v1/assessments", json=body)
    assert resp.status_code == 400
    assert "error" in resp.json()


def test_llm_failure_returns_502(client, fake_llm):
    fake_llm.error = LLMError("llm stopped with max_tokens")
    resp = client.post("/v1/assessments", json=make_request("текст"))
    assert resp.status_code == 502
    # детали — только в логах ai-service (contracts/README.md)
    assert resp.json() == {"error": "service temporarily unavailable"}


def test_urgent_draft_without_recommendations_is_rejected(client, fake_llm):
    fake_llm.draft = LLMDraft(urgency="priority", recommendations=[])
    resp = client.post("/v1/assessments", json=make_request("текст"))
    assert resp.status_code == 502


def test_llm_gets_knowledge_built_from_mis_catalog(client, fake_llm, kb):
    fake_llm.draft = LLMDraft(urgency="normal", recommendations=[rec("SCREENING-ANNUAL")])
    client.post("/v1/assessments", json=make_request("Без патологии"))
    _, _, used_kb = fake_llm.calls[0]
    assert used_kb.catalog_version == kb.catalog_version
    assert set(used_kb.catalog) == set(kb.catalog)


def test_mis_catalog_unavailable_returns_502(client, fake_llm, fake_catalog):
    fake_catalog.catalog = None
    resp = client.post("/v1/assessments", json=make_request("текст"))
    assert resp.status_code == 502
    assert resp.json() == {"error": "service temporarily unavailable"}
    assert fake_llm.calls == []


def test_code_absent_in_mis_catalog_becomes_not_in_clinic(client, fake_llm):
    # например, справочник МИС обновился между ответом модели и постобработкой
    fake_llm.draft = LLMDraft(
        urgency="planned",
        recommendations=[rec("ECHO-CG", service_name="Эхокардиография")],
    )
    recs = client.post("/v1/assessments", json=make_request("Кардиомегалия.")).json()["recommendations"]
    assert [(r["service_code"], r["service_name"]) for r in recs] == [("", "Эхокардиография")]


def test_healthz(client, kb):
    body = client.get("/healthz").json()
    assert body["catalog"] == {"version": kb.catalog_version, "services": len(kb.catalog)}
    assert body["guidelines_version"] == kb.guidelines_version


def test_healthz_before_catalog_loaded(client, fake_catalog):
    fake_catalog.catalog = None
    assert client.get("/healthz").json()["catalog"] is None
