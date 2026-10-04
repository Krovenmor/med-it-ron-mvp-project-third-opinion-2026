# Контракты API

| Документ | Что описывает | Кто потребитель |
|---|---|---|
| [b2b-service.md](b2b-service.md) | main-B2B: приём заключений, карточка кейса, очередь и проверка врачом, план пациента, регистрация записи, демо-часы | МИС, веб врача, b2c-service |
| [b2c-service.md](b2c-service.md) | Граф здоровья пациента, слоты, запись | Веб пациента |
| [mis-demo.md](mis-demo.md) | Демо-МИС: история пациента, слоты, запись, отправка демо-заключений | b2b-service, b2c-service, демо-кнопки |
| [ai-service.md](ai-service.md) | Оценка заключения: срочность и рекомендации | b2b-service |

## Сервисы и связи

| Сервис | Порт | Хранилище |
|---|---|---|
| b2b-service | 8080 | Postgres |
| mis-demo | 8081 | в памяти, сбрасывается при перезапуске |
| b2c-service | 8082 | нет, собирает данные на лету |
| ai-service | задаётся `AI_SERVICE_URL` | — |

```
            POST /api/v1/reports
 mis-demo ───────────────────────────▶ b2b-service ──── POST /v1/assessments ───▶ ai-service
    ▲   ◀── GET /patients/{id}/history ──┘    ▲  ▲
    │                                         │  └──── веб врача
    │   GET history, slots                    │
    │   POST appointments                     │  GET /patients/{id}/plan
    └──────────────────────────── b2c-service ┘  POST /cases/{id}/bookings
                                       ▲
                                       └──── веб пациента
```

## Общие соглашения

- **Формат:** JSON в UTF-8, `Content-Type: application/json`. Неизвестные поля в запросах игнорируются.
- **Дата и время:** RFC 3339 со смещением. b2b-service отдаёт время в UTC (`2026-10-03T09:00:00Z`); слоты и записи МИС идут в часовом поясе клиники (`2026-10-05T09:00:00+03:00`), b2c-service передаёт их как есть. Клиент учитывает смещение и не полагается на конкретный пояс. Дата без времени (дата рождения) — `YYYY-MM-DD`.
- **Идентификаторы:** сущности b2b-service (кейс, рекомендация, запись в b2b) — UUID. Идентификаторы МИС (пациент, исследование, запись, слот) — непрозрачные строки, клиент их не разбирает.
- **Ошибки:** тело всегда `{"error": "описание"}`. Для 400 и 409 описание пригодно для показа разработчику, для 500 и 502 — общее (`internal error`, `service temporarily unavailable`), детали только в логах сервиса.
- **Авторизации в MVP нет.** Действия врача требуют заголовок `X-User-ID` с идентификатором врача: он попадает в аудит. Пациент в B2C идентифицируется ID пациента в МИС в пути запроса — предполагается, что личный кабинет клиники уже авторизовал его.
- **Пациенту не показываются** диагнозы, находки, текст заключения, обоснования для врача и уровень срочности. API, которые смотрят на пациента (план в b2b-service и весь b2c-service), этих полей не содержат.

### Коды ответов

| Код | Когда |
|---|---|
| 200 | Успех; для идемпотентных операций — повтор уже выполненного запроса |
| 201 | Создан новый ресурс |
| 204 | Успех без тела |
| 400 | Некорректный JSON или значение поля; в описании указано поле |
| 401 | Нет заголовка `X-User-ID` у действия врача |
| 404 | Ресурс не найден, в том числе некорректный UUID в пути |
| 409 | Действие недопустимо в текущем состоянии (кейс не на проверке, слот занят, рекомендация уже записана и т. п.) |
| 500 | Внутренняя ошибка |
| 502 | b2c-service: недоступен b2b-service или МИС |

### Идемпотентность

| Операция | Ключ | Повтор |
|---|---|---|
| `POST /api/v1/reports` (b2b) | `source_system` + `study.id` | тот же кейс, 200; другое содержимое исследования — 409 |
| `POST /api/v1/cases/{id}/bookings` (b2b) | `appointment_id` | та же запись, 200 |
| `POST /api/v1/appointments` (mis-demo) | слот + пациент + `referral_id` | та же запись, 200 |
| `POST .../appointments` (b2c) | рекомендация | 409 «уже записаны»; перед ответом существующая запись повторно регистрируется в b2b |

## Справочник значений

| Поле | Значения |
|---|---|
| `modality` | `CT` — КТ, `DX` — рентгенография, `MG` — маммография |
| `sex` | `male`, `female` |
| `urgency` | `normal` — норма, `planned` — плановый, `priority` — приоритетный, `emergency` — неотложный |
| `status` кейса | `draft`, `in_review`, `confirmed`, `notified`, `booked`, `completed`, `declined`, `unreachable` |
| `source` рекомендации | `ai` — от ai-service, `doctor` — добавлена врачом |
| `importance` | `high`, `low` — предварительная метка ai-service |
| `mark` (отметка врача) | `critical` — крайне важна, `minor` — менее важна, `rejected` — не нужна |
| `reject_reason` | `contraindicated` — противопоказано, `other` — другое (обязателен `reject_comment`) |
| `channel` записи | `self` — пациент записался сам, `operator` — записал оператор |

## Жизненный цикл кейса

| Переход | Чем вызван |
|---|---|
| → `draft` | `POST /api/v1/reports` |
| `draft` → `in_review` | фоновая оценка ai-service; при сбое после всех попыток кейс остаётся в `draft` |
| `in_review` → `confirmed` | `POST /api/v1/cases/{id}/confirm` |
| `confirmed` / `notified` → `booked` | `POST /api/v1/cases/{id}/bookings` |
| → `notified`, `completed`, `declined`, `unreachable` | будут реализованы на шаге с таймерами и задачами оператора |

## Сквозной сценарий

1. `POST mis-demo /demo/reports` — mis-demo отправляет демо-заключения в `POST b2b /api/v1/reports`.
2. b2b-service в фоне берёт историю пациента из `GET mis-demo /api/v1/patients/{id}/history` и вызывает `POST ai-service /v1/assessments`; кейс попадает в очередь врача.
3. Веб врача: `GET b2b /api/v1/review/queue` → `GET /api/v1/cases/{id}` → `POST .../open` → `PATCH .../recommendations/{rec_id}` по каждой рекомендации → `POST .../confirm`.
4. Веб пациента: `GET b2c /api/v1/patients/{id}/graph` → `GET .../recommendations/{rec_id}/slots` → `POST .../recommendations/{rec_id}/appointments`.
5. b2c-service записывает пациента в `POST mis-demo /api/v1/appointments` и регистрирует запись в `POST b2b /api/v1/cases/{id}/bookings`; кейс переходит в `booked`.
