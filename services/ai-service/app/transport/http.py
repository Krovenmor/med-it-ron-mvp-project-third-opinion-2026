"""HTTP API по контракту contracts/ai-service.md."""

import logging
from collections.abc import Callable
from contextlib import AbstractAsyncContextManager

from fastapi import FastAPI, Request
from fastapi.exceptions import RequestValidationError
from fastapi.responses import JSONResponse

from app.domain.assessment import AssessmentRequest, AssessmentResponse, ErrorResponse
from app.domain.errors import CatalogUnavailableError, InvalidDraftError, LLMError
from app.service.assessment import AssessmentService

log = logging.getLogger(__name__)

_UNAVAILABLE = "service temporarily unavailable"


def new_app(
    service: AssessmentService,
    health: Callable[[], dict],
    lifespan: Callable[[FastAPI], AbstractAsyncContextManager[None]],
) -> FastAPI:
    app = FastAPI(title="ai-service", version="1", lifespan=lifespan)

    # Ошибки — {"error": "..."}; для 502 и 500 описание общее, детали только в логах
    # (contracts/README.md). Любой не-200 main-B2B считает сбоем и повторяет запрос.
    @app.exception_handler(RequestValidationError)
    async def invalid_request(_: Request, exc: RequestValidationError):
        return _error(400, f"invalid request: {exc.errors()}")

    @app.exception_handler(CatalogUnavailableError)
    async def catalog_unavailable(_: Request, exc: CatalogUnavailableError):
        log.error("mis catalog unavailable: %s", exc)
        return _error(502, _UNAVAILABLE)

    @app.exception_handler(LLMError)
    async def llm_failed(_: Request, exc: LLMError):
        log.error("llm failed: %s", exc)
        return _error(502, _UNAVAILABLE)

    @app.exception_handler(InvalidDraftError)
    async def invalid_draft(_: Request, exc: InvalidDraftError):
        log.error("invalid llm draft: %s", exc)
        return _error(502, _UNAVAILABLE)

    @app.exception_handler(Exception)
    async def internal(_: Request, exc: Exception):
        log.exception("assessment failed")
        return _error(500, "internal error")

    @app.post(
        "/v1/assessments",
        response_model=AssessmentResponse,
        responses={400: {"model": ErrorResponse}, 502: {"model": ErrorResponse}},
    )
    async def assess(req: AssessmentRequest) -> AssessmentResponse:
        return await service.assess(req)

    @app.get("/healthz")
    async def healthz() -> dict:
        return health()

    return app


def _error(status: int, message: str) -> JSONResponse:
    return JSONResponse(status_code=status, content={"error": message})
