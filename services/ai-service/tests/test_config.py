import os

import pytest

from app import config
from app.config import load_env_files, load_settings


def test_settings_from_default_env():
    s = load_settings()
    assert (s.http_host, s.http_port) == ("0.0.0.0", 8083)
    assert s.llm_model == "GigaChat-2-Max"
    assert s.mis_url == "http://localhost:8081"
    assert s.knowledge_dir == config.BASE_DIR / "knowledge"


def test_process_env_wins_over_files(monkeypatch):
    monkeypatch.setenv("MIS_URL", "http://mis-demo:8081")
    monkeypatch.setenv("HTTP_ADDR", "127.0.0.1:9000")
    s = load_settings()
    assert s.mis_url == "http://mis-demo:8081"
    assert (s.http_host, s.http_port) == ("127.0.0.1", 9000)


def test_first_file_wins_and_empty_values_are_skipped(tmp_path, monkeypatch):
    for key in ("AI_TEST_SECRET", "AI_TEST_EMPTY", "AI_TEST_SHARED"):
        monkeypatch.delenv(key, raising=False)
    local = tmp_path / ".env"
    local.write_text("AI_TEST_SECRET=from-dotenv\nAI_TEST_EMPTY=\nAI_TEST_SHARED=from-dotenv\n")
    defaults = tmp_path / "default.env"
    defaults.write_text("AI_TEST_EMPTY=from-default\nAI_TEST_SHARED=from-default\n")

    load_env_files((local, defaults, tmp_path / "missing.env"))

    assert os.environ["AI_TEST_SECRET"] == "from-dotenv"
    assert os.environ["AI_TEST_SHARED"] == "from-dotenv"
    # пустое значение в .env не затирает значение из default.env
    assert os.environ["AI_TEST_EMPTY"] == "from-default"
    for key in ("AI_TEST_SECRET", "AI_TEST_EMPTY", "AI_TEST_SHARED"):
        monkeypatch.delenv(key)


def test_missing_settings_are_listed(monkeypatch):
    # значения нет ни в окружении, ни в файлах
    monkeypatch.setattr(config, "ENV_FILES", ())
    monkeypatch.setenv("MIS_URL", " ")
    monkeypatch.setenv("LLM_MODEL", "")
    with pytest.raises(ValueError, match="LLM_MODEL, MIS_URL"):
        load_settings()


def test_empty_process_env_is_filled_from_default_env(monkeypatch):
    # compose передаёт пустую строку `LLM_MODEL=` из .env как пустую переменную окружения
    monkeypatch.setenv("LLM_MODEL", "")
    monkeypatch.setenv("LLM_MAX_TOKENS", "  ")
    s = load_settings()
    assert s.llm_model == "GigaChat-2-Max"
    assert s.llm_max_tokens == 4096


def test_ca_bundle_not_set(monkeypatch):
    monkeypatch.setenv("GIGACHAT_CA_BUNDLE_FILE", "")
    assert load_settings().gigachat_ca_bundle_file is None


def test_ca_bundle_relative_to_service_dir(monkeypatch, tmp_path):
    monkeypatch.setattr(config, "BASE_DIR", tmp_path)
    (tmp_path / "certs").mkdir()
    (tmp_path / "certs" / "ca.pem").write_text("cert")
    monkeypatch.setenv("GIGACHAT_CA_BUNDLE_FILE", "certs/ca.pem")
    assert config._ca_bundle_file() == tmp_path / "certs" / "ca.pem"


def test_ca_bundle_missing_file_fails_fast(monkeypatch):
    monkeypatch.setenv("GIGACHAT_CA_BUNDLE_FILE", "certs/missing.pem")
    with pytest.raises(ValueError, match="GIGACHAT_CA_BUNDLE_FILE"):
        load_settings()


def test_output_mode_is_validated(monkeypatch):
    monkeypatch.setenv("LLM_OUTPUT_MODE", "xml")
    with pytest.raises(ValueError, match="LLM_OUTPUT_MODE"):
        load_settings()
