"""Настройки сервиса. Как в остальных сервисах проекта, значений по умолчанию в коде нет.

Откуда берутся значения (выше — важнее):
1. переменные окружения процесса (в docker-compose — из env_file);
2. .env — секреты и настройки конкретной машины, в .gitignore;
3. config/default.env — несекретные настройки, коммитится и лежит в образе.
Пустое значение — и в файле, и в окружении — считается незаданным: строка `LLM_MODEL=`
в .env не перекрывает значение из default.env, даже когда compose передал её в контейнер.
"""

import os
from dataclasses import dataclass
from pathlib import Path

from dotenv import dotenv_values

BASE_DIR = Path(__file__).resolve().parent.parent
ENV_FILES: tuple[Path, ...] = (BASE_DIR / ".env", BASE_DIR / "config" / "default.env")


@dataclass(frozen=True)
class Settings:
    http_host: str
    http_port: int
    llm_model: str
    llm_timeout: float
    llm_max_tokens: int
    # Сколько запросов к LLM выполняется одновременно; 0 — без ограничения.
    llm_max_concurrency: int
    llm_temperature: float
    # 1.0 — без штрафа за повторы: в JSON ключи повторяются в каждой рекомендации,
    # и штраф заставляет модель их коверкать.
    llm_repetition_penalty: float
    # Как получать от модели JSON: prompt — схема описана в промпте, ответ обычным текстом;
    # json_schema — режим response_format GigaChat (beta, на практике портит JSON).
    llm_output_mode: str
    # Справочник услуг клиники — из МИС (GET {mis_url}/api/v1/services).
    mis_url: str
    mis_timeout: float
    # Сколько секунд справочник считается свежим, потом запрашивается заново.
    catalog_ttl: float
    knowledge_dir: Path
    log_level: str
    # Корневой сертификат НУЦ Минцифры для GigaChat; None — системные сертификаты.
    gigachat_ca_bundle_file: Path | None


def load_env_files(files: tuple[Path, ...] | None = None) -> None:
    """Дополняет окружение значениями из файлов, не перетирая уже заданные непустые."""
    for path in ENV_FILES if files is None else files:
        if not path.is_file():
            continue
        for key, value in dotenv_values(path).items():
            if value and not os.environ.get(key, "").strip():
                os.environ[key] = value


def load_settings() -> Settings:
    load_env_files()
    env = _Env()
    host, port = _parse_addr(env.str("HTTP_ADDR"))
    knowledge_dir = Path(env.str("KNOWLEDGE_DIR"))
    settings = Settings(
        http_host=host,
        http_port=port,
        llm_model=env.str("LLM_MODEL"),
        llm_timeout=env.float("LLM_TIMEOUT"),
        llm_max_tokens=env.int("LLM_MAX_TOKENS"),
        llm_max_concurrency=env.int("LLM_MAX_CONCURRENCY"),
        llm_temperature=env.float("LLM_TEMPERATURE"),
        llm_repetition_penalty=env.float("LLM_REPETITION_PENALTY"),
        llm_output_mode=env.choice("LLM_OUTPUT_MODE", ("prompt", "json_schema")),
        mis_url=env.str("MIS_URL"),
        mis_timeout=env.float("MIS_TIMEOUT"),
        catalog_ttl=env.float("CATALOG_TTL"),
        knowledge_dir=knowledge_dir if knowledge_dir.is_absolute() else BASE_DIR / knowledge_dir,
        log_level=env.str("LOG_LEVEL"),
        gigachat_ca_bundle_file=_ca_bundle_file(),
    )
    env.raise_if_missing()
    return settings


def _ca_bundle_file() -> Path | None:
    """GIGACHAT_CA_BUNDLE_FILE: абсолютный путь или относительно папки сервиса (certs/...).

    Относительный путь одинаково работает локально и в контейнере, куда certs/ монтируется томом.
    """
    value = os.environ.get("GIGACHAT_CA_BUNDLE_FILE", "").strip()
    if not value:
        return None
    path = Path(value) if Path(value).is_absolute() else BASE_DIR / value
    if not path.is_file():
        # SDK при заданном пути игнорирует GIGACHAT_VERIFY_SSL_CERTS — лучше упасть сразу и понятно.
        raise ValueError(f"GIGACHAT_CA_BUNDLE_FILE: file not found: {path}")
    return path


class _Env:
    """Читает обязательные переменные и копит список отсутствующих, чтобы назвать все сразу."""

    def __init__(self) -> None:
        self.missing: list[str] = []

    def str(self, key: str) -> str:
        value = os.environ.get(key, "").strip()
        if not value:
            self.missing.append(key)
        return value

    def int(self, key: str) -> int:
        value = self.str(key)
        return int(value) if value else 0

    def float(self, key: str) -> float:
        value = self.str(key)
        return float(value) if value else 0.0

    def choice(self, key: str, allowed: tuple[str, ...]) -> str:
        value = self.str(key)
        if value and value not in allowed:
            raise ValueError(f"{key} must be one of {', '.join(allowed)}, got {value!r}")
        return value

    def raise_if_missing(self) -> None:
        if self.missing:
            names = ", ".join(self.missing)
            raise ValueError(f"required settings are not set: {names} (see config/default.env)")


def _parse_addr(addr: str) -> tuple[str, int]:
    """HTTP_ADDR в формате остальных сервисов: ":8083" или "127.0.0.1:8083"."""
    if not addr:
        return "", 0
    host, _, port = addr.rpartition(":")
    return host or "0.0.0.0", int(port)
