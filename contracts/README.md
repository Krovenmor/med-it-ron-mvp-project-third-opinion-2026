# Контракты API

| Документ | Что описывает | Кто потребитель |
|---|---|---|
| [b2b-service.md](b2b-service.md) | main-B2B: приём заключений, карточка кейса, очередь и проверка врачом, протокол таймеров, задачи администратора, журнал уведомлений, дашборд руководителя, план пациента, регистрация записи, демо-часы и демо-история | МИС, веб врача, администратора и руководителя, b2c-service |
| [b2c-service.md](b2c-service.md) | Маршрут пациента, слоты, запись, отказ, просьба перезвонить | Веб пациента |
| [mis-demo.md](mis-demo.md) | Демо-МИС: история пациента, слоты, запись, отправка демо-заключений | b2b-service, b2c-service, демо-кнопки |
| [ai-service.md](ai-service.md) | Оценка заключения: срочность и рекомендации | b2b-service |

## Сервисы и связи

| Сервис | Порт | Хранилище |
|---|---|---|
| b2b-service | 8080 | Postgres |
| mis-demo | 8081 | в памяти, сбрасывается при перезапуске |
| b2c-service | 8082 | нет, собирает данные на лету |
| ai-service | 8083 | нет; справочник услуг берёт из mis-demo |
| frontend-service | 3000 | нет; nginx отдаёт статику и проксирует `/b2b/`, `/mis/`, `/b2c/` на сервисы |

```
            POST /api/v1/reports
 mis-demo ───────────────────────────▶ b2b-service ──── POST /v1/assessments ───▶ ai-service
    ▲   ◀── GET history, slots ──────────┘    ▲  ▲
    │       POST appointments (администратор) │  └──── веб врача, администратора, руководителя
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
- **Авторизации в MVP нет.** Действия врача и администратора требуют заголовок `X-User-ID` с идентификатором сотрудника: он попадает в аудит. Пациент в B2C идентифицируется ID пациента в МИС в пути запроса — предполагается, что личный кабинет клиники уже авторизовал его.
- **Пациенту не показываются** диагнозы, находки, текст заключения, обоснования для врача и уровень срочности. API, которые смотрят на пациента (план в b2b-service и весь b2c-service), этих полей не содержат.

### Коды ответов

| Код | Когда |
|---|---|
| 200 | Успех; для идемпотентных операций — повтор уже выполненного запроса |
| 201 | Создан новый ресурс |
| 204 | Успех без тела |
| 400 | Некорректный JSON или значение поля; в описании указано поле |
| 401 | Нет заголовка `X-User-ID` у действия врача или администратора |
| 404 | Ресурс не найден, в том числе некорректный UUID в пути |
| 409 | Действие недопустимо в текущем состоянии (кейс не на проверке, слот занят, рекомендация уже записана и т. п.) |
| 500 | Внутренняя ошибка |
| 502 | Недоступен вызываемый сервис: МИС для записи администратором в b2b-service; b2b-service или МИС для b2c-service |

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
| `status` заключения (gateway) | `received`, `assessed` |
| `status` в карточке врача (doctor) | `in_review`, `confirmed` |
| `status` кейса в плане и у администратора (care) | `received`, `in_review`, `confirmed`, `notified`, `booked`, `completed`, `declined`, `unreachable` |
| `source` рекомендации | `ai` — от ai-service, `doctor` — добавлена врачом |
| `importance` | `high`, `low` — предварительная метка ai-service |
| `mark` (отметка врача) | `critical` — крайне важна, `minor` — менее важна, `rejected` — не нужна |
| `reject_reason` | `contraindicated` — противопоказано, `other` — другое (обязателен `reject_comment`) |
| `channel` записи | `self` — пациент записался сам, `operator` — записал администратор |
| `reason` задачи | `emergency` — неотложный случай, `no_booking` — не записался вовремя, `help_request` — пациент просит перезвонить |
| `status` задачи | активные: `new` — новая, `in_progress` — в работе, `no_answer` — не дозвонились, `callback` — перезвонить; закрытые: `contacted` — дозвонились, `booked` — записан, `declined` — отказ, `handed_to_doctor` — передано врачу, `cancelled` — снята протоколом |
| `outcome` попытки | `no_answer`, `callback`, `contacted`, `declined`, `booked`, `handed_to_doctor` |
| `decline_reason` | `expensive` — дорого, `far` — далеко, `other_clinic` — лечится в другой клинике, `not_needed` — не считает нужным, `other` — другое (обязателен `comment`) |
| `recipient` уведомления | `patient`, `duty_doctor` — дежурный врач, `admin_on_duty` — дежурный администратор |
| `channel` уведомления | `push`, `sms`, `staff` — внутреннее сотрудникам |
| `kind` уведомления | `plan_ready`, `no_findings`, `reminder`, `urgent_contact`, `unreachable`, `review_escalation`, `contact_escalation`, `handed_to_doctor` |

## Жизненный цикл кейса

| Переход | Чем вызван |
|---|---|
| → `received` | `POST /api/v1/reports` |
| `received` → `in_review` | фоновая оценка ai-service; при сбое после всех попыток заключение остаётся `received`, врач кейс не видит |
| `in_review` → `confirmed` | `POST /api/v1/cases/{id}/confirm` |
| `confirmed` → `notified` | фоновое уведомление пациенту сразу после подтверждения |
| `confirmed` / `notified` → `unreachable` | третья попытка администратора без ответа |
| `confirmed` / `notified` / `unreachable` → `declined` | отказ пациента, записанный администратором |
| `confirmed` / `notified` / `unreachable` / `declined` → `booked` | запись пациента (`POST /api/v1/cases/{id}/bookings`) или администратора (`POST /api/v1/operator/tasks/{id}/bookings`) |
| `booked` → `completed` | не реализован в MVP: нужен факт визита из МИС; встречается только в демо-истории |

## Сквозной сценарий

1. `POST mis-demo /demo/reports` — mis-demo отправляет демо-заключения в `POST b2b /api/v1/reports`.
2. b2b-service в фоне берёт историю пациента из `GET mis-demo /api/v1/patients/{id}/history` и вызывает `POST ai-service /v1/assessments`; кейс попадает в очередь врача.
3. Веб врача: `GET b2b /api/v1/review/queue` → `GET /api/v1/cases/{id}` → `POST .../open` → `PATCH .../recommendations/{rec_id}` по каждой рекомендации → `POST .../confirm`.
4. Веб пациента: `GET b2c /api/v1/patients/{id}/route` → `GET .../recommendations/{rec_id}/slots` → `POST .../recommendations/{rec_id}/appointments` (или `.../decline`, `.../help-requests`).
5. b2c-service записывает пациента в `POST mis-demo /api/v1/appointments` и регистрирует запись в `POST b2b /api/v1/cases/{id}/bookings`; кейс переходит в `booked`.
6. Если пациент не записался, протокол создаёт задачу администратору. Веб администратора: `GET b2b /api/v1/operator/tasks` → `GET .../tasks/{id}` → `POST .../take` → результат звонка или `POST .../bookings`; b2b-service сам записывает пациента в `POST mis-demo /api/v1/appointments`. Журнал — `GET b2b /api/v1/notifications`.
