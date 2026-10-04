# ai-service

Формирует рекомендации маршрута пациента по заключению ИИ. Контракт: [contracts/ai-service.md](../../contracts/ai-service.md).

## Структура

Слои — как в Go-сервисах проекта (`internal/domain`, `service`, `infra`, `transport`, `di`):

```
ai-service/
├── app/
│   ├── __main__.py         точка входа: python -m app
│   ├── main.py             сборка зависимостей (аналог internal/di)
│   ├── config.py           настройки из окружения и env-файлов
│   ├── domain/             модели и правила, без ввода-вывода
│   │   ├── assessment.py     запрос/ответ по контракту, черновик LLM
│   │   ├── catalog.py        справочник услуг
│   │   ├── guidelines.py     фрагменты клинических рекомендаций
│   │   ├── rules.py          маркеры из заключения и нижняя граница срочности
│   │   └── errors.py
│   ├── service/            сценарий оценки
│   │   ├── ports.py          интерфейсы внешних зависимостей: LLM, CatalogSource
│   │   ├── assessment.py     база знаний → правила → LLM → постобработка
│   │   ├── knowledge.py      промпт и JSON-схема ответа из справочника и КР
│   │   └── postprocess.py    проверки ответа модели
│   ├── infra/              адаптеры
│   │   ├── mis.py            справочник услуг из МИС (GET /api/v1/services)
│   │   ├── gigachat.py       LLM на GigaChat
│   │   ├── llm.py            общее для адаптеров LLM
│   │   └── guidelines.py     чтение knowledge/
│   └── transport/
│       └── http.py           FastAPI: POST /v1/assessments, GET /healthz
├── knowledge/              клинические рекомендации и инструкция для модели (данные, не код)
├── config/
│   ├── default.env         несекретные настройки для локального запуска
│   └── docker.env          переопределения для docker-compose
├── certs/                  сертификат НУЦ Минцифры для GigaChat (монтируется в контейнер)
├── .env.example            шаблон секретов → копируется в .env
├── tests/
├── Dockerfile
└── pyproject.toml
```

Зависимости направлены внутрь: `transport → service → domain`, `infra` реализует порты из `service/ports.py`.
Другую модель вместо GigaChat можно подключить новым адаптером в `infra/` и одной строкой в `main.py`.

## Как работает оценка

1. **Справочник услуг** — из МИС ([contracts/mis-demo.md](../../contracts/mis-demo.md)): код, название, описание для пациента.
   Своей копии нет: рекомендовать можно только то, на что МИС записывает. Кэш на `CATALOG_TTL` секунд;
   не удалось обновить — работаем с последней версией; МИС недоступна с самого старта — оценка отвечает 502,
   и main-B2B повторит запрос. `catalog_version` = `mis-<хеш справочника>`.
2. **Правила** (`domain/rules.py`) находят в заключении BI-RADS, размер очага, находки и задают нижнюю границу срочности.
3. **LLM** получает заключение, справочник и выдержки КР и отвечает JSON по схеме, где коды услуг и ссылки на КР
   ограничены `enum`.
4. **Постобработка** берёт названия из справочника, проставляет `already_booked`, убирает дубли и уже рекомендованное,
   не даёт опустить срочность ниже правил и заменяет `patient_text` с диагнозами.

**Выдержки КР в `knowledge/guidelines/` и пороги в `domain/rules.py` — черновик.** Их должен сверить медицинский эксперт.

## Настройки: config/default.env, .env и .env.example

Как в остальных сервисах проекта, значений по умолчанию в коде нет — всё в env-файлах:

| Файл | Что в нём | В git |
|---|---|---|
| `config/default.env` | все несекретные настройки для локального запуска | да |
| `config/docker.env` | что меняется в docker-compose (адрес МИС в docker-сети) | да |
| `.env.example` | **шаблон** секретов: имена переменных без значений | да |
| `.env` | **настоящие** секреты и настройки вашей машины (ключ GigaChat, путь к сертификату) | нет, в `.gitignore` |

Приоритет: переменные окружения → `.env` → `config/default.env`. Пустое значение в файле считается незаданным.
Любую настройку из `default.env` можно переопределить у себя в `.env`, не трогая общий файл.

## Запуск

Нужна запущенная mis-demo (`MIS_URL`, по умолчанию `http://localhost:8081`).

```bash
python -m venv .venv && .venv/bin/pip install -e ".[dev]"
cp .env.example .env      # вписать GIGACHAT_CREDENTIALS и GIGACHAT_CA_BUNDLE_FILE
.venv/bin/python -m app   # слушает HTTP_ADDR=:8083
```

В b2b-service: `AI_SERVICE_MOCK=false`, `AI_SERVICE_URL=http://localhost:8083`.
Проверка: `curl localhost:8083/healthz` — модель, загруженный справочник (`null`, пока МИС не ответила), версия КР.

Ручная проверка — готовые тела запросов в `examples/` (демо-пациенты mis-demo и заключения в стиле Kafka):

```bash
curl -s -X POST localhost:8083/v1/assessments -H 'Content-Type: application/json' \
  -d @examples/01-anna-mg-birads4.json | jq
```

Тесты (GigaChat и МИС подменяются фейками; локальный `.env` тесты не читают):

```bash
.venv/bin/pytest
```

### GigaChat

1. developers.sber.ru → Studio → создать проект **GigaChat API** → «Получить ключ».
   Это Authorization key, он показывается один раз → `GIGACHAT_CREDENTIALS` в `.env`.
2. Физлицам доступен бесплатный пакет токенов (Freemium); объёмы и модели — в личном кабинете.
   Один кейс — порядка 4–6 тыс. токенов. Параллельность на бесплатном тарифе ограничена, поэтому
   `LLM_MAX_CONCURRENCY=1`: запросы от воркеров main-B2B выстроятся в очередь, а не получат 429.
3. Скачать корневой сертификат НУЦ Минцифры (gosuslugi.ru/crt), положить в `certs/` и указать в `.env`
   `GIGACHAT_CA_BUNDLE_FILE=certs/russian_trusted_root_ca.pem` — путь относительно папки сервиса, работает
   и локально, и в Docker. Без сертификата — `CERTIFICATE_VERIFY_FAILED`; временно, только локально,
   можно `GIGACHAT_VERIFY_SSL_CERTS=False`.
4. Список доступных моделей: `GigaChat().get_models()`; нужную — в `LLM_MODEL`.

Режим `response_format` (JSON-схема) у GigaChat помечен как beta и на практике вставляет мусор в JSON,
поэтому по умолчанию схема передаётся текстом в промпте (`LLM_OUTPUT_MODE=prompt`). Если ответ не разобрался,
адаптер один раз переспрашивает модель; не помогло — 502, и main-B2B повторит запрос. Постобработка отбрасывает
неизвестные коды услуг и ссылки на КР.

### Docker

ai-service есть в корневом `docker-compose.yml`, порт 8083 открыт только на localhost. b2b-service обращается к нему
(`AI_SERVICE_MOCK=false`, `AI_SERVICE_URL=http://ai-service:8083` в `services/b2b-service/config/docker.env`).

```bash
docker compose up --build        # из корня репозитория
```

Настройки контейнер получает из `env_file` (позже — важнее): `config/default.env` → `.env` → `config/docker.env`.
`.env` необязателен: без него сервис стартует, но оценки будут падать, кейсы останутся в `draft`.
Пустые строки вида `LLM_MODEL=` в `.env` не перекрывают `default.env`. Папка `certs/` монтируется в `/srv/certs`.

Вернуть мок ai-service в b2b (если нет ключа GigaChat): `AI_SERVICE_MOCK=true` в
`services/b2b-service/config/docker.env`.

## Переменные

| Переменная | Где задана | |
|---|---|---|
| `HTTP_ADDR` | default.env: `:8083` | адрес HTTP-сервера |
| `LLM_MODEL` | default.env: `GigaChat-2-Max` | |
| `LLM_TIMEOUT` | default.env: `80` | секунд; main-B2B ждёт 90 с |
| `LLM_MAX_TOKENS` | default.env: `4096` | |
| `LLM_MAX_CONCURRENCY` | default.env: `1` | одновременных запросов к LLM, `0` — без ограничения |
| `LLM_TEMPERATURE` | default.env: `0.1` | |
| `LLM_REPETITION_PENALTY` | default.env: `1.0` | штраф за повторы |
| `LLM_OUTPUT_MODE` | default.env: `prompt` | `prompt` — схема ответа в промпте; `json_schema` — beta-режим GigaChat, портит JSON |
| `GIGACHAT_SCOPE` | default.env: `GIGACHAT_API_PERS` | читает SDK gigachat |
| `GIGACHAT_CREDENTIALS` | .env | Authorization key, читает SDK gigachat |
| `GIGACHAT_CA_BUNDLE_FILE` | .env | сертификат Минцифры: абсолютный путь или относительно папки сервиса |
| `MIS_URL` | default.env, docker.env | справочник услуг |
| `MIS_TIMEOUT` | default.env: `5` | секунд |
| `CATALOG_TTL` | default.env: `300` | секунд |
| `KNOWLEDGE_DIR` | default.env: `knowledge` | относительно папки сервиса |
| `LOG_LEVEL` | default.env: `INFO` | |
