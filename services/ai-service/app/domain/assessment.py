"""Модели запроса и ответа по контракту main-B2B ↔ ai-service (contracts/ai-service.md, v1).

Неизвестные поля во входящем JSON игнорируются (поведение Pydantic по умолчанию),
как того требует контракт.
"""

from datetime import datetime
from typing import Literal

from pydantic import BaseModel

Modality = Literal["CT", "DX", "MG"]
Sex = Literal["male", "female"]
Urgency = Literal["normal", "planned", "priority", "emergency"]
Importance = Literal["high", "low"]

URGENCY_ORDER: tuple[Urgency, ...] = ("normal", "planned", "priority", "emergency")


# --- Запрос от main-B2B ---------------------------------------------------------------


class Study(BaseModel):
    modality: Modality
    body_site: str
    performed_at: datetime


class Patient(BaseModel):
    age: int
    sex: Sex


class Service(BaseModel):
    code: str
    name: str


class Visit(BaseModel):
    service_code: str
    service_name: str
    visited_at: datetime


class Appointment(BaseModel):
    service_code: str
    service_name: str
    scheduled_at: datetime


class AssessmentRequest(BaseModel):
    case_id: str
    study: Study
    conclusion: str
    patient: Patient
    current_recommendations: list[Service]
    visits: list[Visit]
    appointments: list[Appointment]


# --- Ответ main-B2B -------------------------------------------------------------------


class Recommendation(BaseModel):
    service_code: str
    service_name: str
    importance: Importance
    rationale: str
    guideline_ref: str
    patient_text: str
    already_booked: bool


class AssessmentResponse(BaseModel):
    urgency: Urgency
    catalog_version: str
    guidelines_version: str
    recommendations: list[Recommendation]


class ErrorResponse(BaseModel):
    error: str


# --- Черновик от LLM (до постобработки) -----------------------------------------------
# Отличается от ответа: guideline_ref здесь — ID фрагмента из базы КР, а не текст ссылки;
# already_booked и версии LLM не заполняет, их проставляет код.


class DraftRecommendation(BaseModel):
    service_code: str
    service_name: str
    importance: Importance
    rationale: str
    guideline_ref: str
    patient_text: str


class LLMDraft(BaseModel):
    urgency: Urgency
    recommendations: list[DraftRecommendation]
