# Pet-проект: Pulse — сервис мониторинга доступности эндпоинтов

## 1. Идея в одном абзаце

Сервис, в который пользователь добавляет список URL'ов (свои сайты, API, вебхуки), настраивает интервал
и таймаут проверки, а Pulse регулярно «пингует» их по HTTP, сохраняет историю, считает uptime/latency
и отправляет уведомления, когда сервис упал или поднялся обратно. По смыслу — упрощённый Uptime Kuma /
Better Stack / Healthchecks.io.

Почему именно это:

- **Io-bound по природе**: основная нагрузка — параллельные HTTP-запросы наружу, естественный повод
  научиться правильно работать с горутинами, воркер-пулами, контекстами, таймаутами и graceful shutdown.
- **Реально полезно**: такой сервис можно оставить крутиться у себя на VPS и мониторить свои пет-проекты.
- **Законченный продукт**: есть чёткий MVP, который можно «объявить готовым», не увязая в бесконечном скоупе.
- **Покрывает почти весь типовой бэкенд-стек**: REST API, аутентификация, миграции БД, кэш, фоновые задачи,
  внешние интеграции (Telegram/SMTP), метрики, логирование, тесты, Docker, docker-compose.

---

## 2. Стек и инструменты

Обязательная часть (MVP):

- **Язык**: Go 1.22+
- **HTTP-фреймворк**: `net/http` + [chi](https://github.com/go-chi/chi) (или `gin` — на вкус).
  Рекомендую `chi`: он ближе к стандартной библиотеке и учит работать с `http.Handler`.
- **БД**: PostgreSQL 16
- **Драйвер/доступ**: `pgx/v5` + `sqlc` или `pgx` + `squirrel`. Для первого раза можно просто `pgx`
  и руками писать SQL — это полезнее, чем прятаться за ORM.
- **Миграции**: `goose` или `golang-migrate`
- **Кэш / rate limiting**: Redis 7
- **Конфиг**: `envconfig` или `viper` (читаем из env, 12-factor)
- **Логирование**: стандартный `log/slog` (structured JSON-логи)
- **Метрики**: `prometheus/client_golang` (экспортируем `/metrics`)
- **Тесты**: стандартный `testing` + `testify/require`. Для интеграционных — `testcontainers-go`.
- **Контейнеризация**: Docker (multi-stage build, distroless или alpine) + docker-compose

Опционально (на «вторую итерацию»):

- **Swagger/OpenAPI**: `swaggo/swag` или написать `openapi.yaml` руками
- **Миграция на gRPC** одного из эндпоинтов — чтобы пощупать
- **Graceful shutdown через `errgroup`** (не опционально, если честно — лучше сразу)
- **CI**: GitHub Actions с `go test ./...`, `golangci-lint`, сборкой образа

---

## 3. Функциональные требования (MVP)

### 3.1. Пользователи и аутентификация

- Регистрация по email + паролю (пароль хэшируется bcrypt/argon2).
- Логин, выдача **JWT access-токена** (TTL ~15 минут) и **refresh-токена** (TTL ~30 дней, хранится в БД).
- Эндпоинт `POST /auth/refresh` для обновления access-токена.
- Эндпоинт `POST /auth/logout` (инвалидация refresh-токена).
- Вся работа с мониторами — только под авторизованным пользователем.

### 3.2. Мониторы (monitors)

CRUD над сущностью «монитор»:

- `name` — человекочитаемое имя
- `url` — проверяемый URL
- `method` — `GET` / `HEAD` / `POST` (для MVP хватит `GET` и `HEAD`)
- `interval_seconds` — интервал между проверками (например, 30–3600)
- `timeout_seconds` — таймаут запроса (1–30)
- `expected_status` — ожидаемый HTTP-статус (по умолчанию 200)
- `is_active` — вкл/выкл
- `created_at`, `updated_at`

Ограничения: один пользователь видит и трогает только свои мониторы (проверяется в middleware/хэндлере).

### 3.3. Проверки (checks) — история

Каждая проба сохраняется как запись:

- `monitor_id`
- `checked_at`
- `status` — `up` / `down`
- `http_status` — фактический код ответа (nullable, если запрос вообще не дошёл)
- `latency_ms` — время ответа
- `error` — текст ошибки (таймаут, DNS, TLS, non-2xx, ...), nullable

Эндпоинты:

- `GET /monitors/{id}/checks?limit=100&from=...&to=...` — постранично история.
- `GET /monitors/{id}/stats?period=24h` — агрегаты: uptime %, avg latency, p95 latency, количество инцидентов.

### 3.4. Инциденты (incidents)

Инцидент открывается при первой `down`-пробе, закрывается при первой `up`-пробе после этого.
Хранится как интервал `[started_at, ended_at]`. Нужен, чтобы не слать 500 уведомлений, пока сервис лежит.

- `GET /monitors/{id}/incidents`
- Отдельной ручки создания не нужно — создаётся воркером.

### 3.5. Уведомления (notifications)

Каналы (MVP — достаточно одного, например Telegram):

- **Telegram bot**: пользователь добавляет свой `chat_id` и токен бота (или общий бот сервиса).
- **Webhook**: `POST` на произвольный URL с JSON-пейлоадом.
- **(Опционально) SMTP** — отправка на email.

Правила доставки:

- Уведомление «упал» отправляется при открытии инцидента, **но не раньше чем через N сек подряд вниз**
  (чтобы не реагировать на одну флаковую пробу). Типично: down-подтверждение через 2–3 failed пробы подряд.
- Уведомление «поднялся» — при закрытии инцидента.
- В уведомлении должна быть инфа: имя монитора, URL, когда упал, как долго лежит, последняя ошибка.

### 3.6. Фоновый рабочий процесс (scheduler / worker)

Это сердце сервиса и главное, ради чего проект стоит делать.

- При старте приложения поднимается **один планировщик**, который раз в X секунд тянет из БД активные
  мониторы, у которых `next_check_at <= now()`.
- Планировщик кладёт задачи в **канал / очередь**. Её разгребает **пул воркеров** (N штук, конфигурируемо).
- Каждый воркер делает HTTP-запрос с таймаутом (контекст!), пишет результат пробы в БД, обновляет
  `next_check_at`, при необходимости открывает/закрывает инцидент и ставит задачу на отправку уведомления.
- Все горутины слушают `context.Context` главного приложения → graceful shutdown по SIGINT/SIGTERM.

Важные нюансы, которые полезно прочувствовать:

- Один монитор не должен «залипать» и блокировать весь пул (таймаут обязателен).
- Если пользователь поставил 1000 мониторов с интервалом 30 сек, пул должен справляться — добавить
  ограничение по количеству одновременно выполняемых проб (семафор/`chan struct{}`).
- На уровне HTTP-клиента — переиспользуемый `http.Client` с настроенным `Transport` (keep-alive,
  лимит коннектов). Именно тут io-bound раскрывается на полную.

### 3.7. Health & metrics

- `GET /healthz` — живой ли процесс (всегда 200, если процесс жив).
- `GET /readyz` — готов ли принимать трафик (есть коннект к БД/Redis).
- `GET /metrics` — Prometheus-формат:
  - `pulse_checks_total{status="up|down"}`
  - `pulse_check_duration_seconds` (histogram)
  - `pulse_active_monitors`
  - `pulse_notifications_sent_total{channel="telegram|webhook"}`
  - `pulse_worker_queue_depth`

---

## 4. Архитектура и структура кода

Рекомендуемая раскладка (можно адаптировать, но стоит следовать идее «разделение слоёв»):

```
pulse/
├── cmd/
│   └── pulse/
│       └── main.go              # точка входа, сборка зависимостей (composition root)
├── internal/
│   ├── config/                  # загрузка env
│   ├── http/
│   │   ├── handlers/            # хэндлеры REST API
│   │   ├── middleware/          # auth, logging, recover, request-id
│   │   └── router.go
│   ├── auth/                    # JWT, хэширование паролей
│   ├── monitor/                 # доменная логика мониторов (сервис-слой)
│   ├── checker/                 # пул воркеров, HTTP-пробы
│   ├── scheduler/               # планировщик, который тикает и ставит задачи
│   ├── notifier/                # интерфейс Notifier + реализации (telegram, webhook)
│   ├── storage/
│   │   ├── postgres/            # репозитории (pgx), миграции
│   │   └── redis/               # кэш, rate-limit
│   └── observability/           # slog, prometheus
├── migrations/                  # *.sql для goose/migrate
├── deployments/
│   ├── docker-compose.yml
│   └── Dockerfile
├── api/
│   └── openapi.yaml             # (опционально) спецификация API
├── Makefile                     # run, test, lint, migrate, docker-up
├── go.mod / go.sum
└── README.md
```

Принципы:

- **Handlers** тонкие: парсинг запроса → вызов сервиса → рендер ответа.
- **Сервисный слой** (`monitor`, `auth`, …) ничего не знает про HTTP и БД — работает с интерфейсами
  репозиториев. Это сразу даёт тестируемость.
- **Репозитории** в `storage/postgres` — реализуют интерфейсы, живущие в доменных пакетах.
- **`cmd/pulse/main.go`** собирает всё руками (никаких DI-фреймворков). Это важная практика для Go.

---

## 5. Схема БД (стартовая, SQL-псевдокод)

```sql
CREATE TABLE users (
    id          BIGSERIAL PRIMARY KEY,
    email       TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE refresh_tokens (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ
);

CREATE TABLE monitors (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    url             TEXT NOT NULL,
    method          TEXT NOT NULL DEFAULT 'GET',
    interval_seconds INT NOT NULL,
    timeout_seconds  INT NOT NULL,
    expected_status  INT NOT NULL DEFAULT 200,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    next_check_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    consecutive_failures INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_monitors_due ON monitors (next_check_at) WHERE is_active;

CREATE TABLE checks (
    id          BIGSERIAL PRIMARY KEY,
    monitor_id  BIGINT NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    checked_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    status      TEXT NOT NULL,            -- 'up' | 'down'
    http_status INT,
    latency_ms  INT,
    error       TEXT
);
CREATE INDEX idx_checks_monitor_time ON checks (monitor_id, checked_at DESC);

CREATE TABLE incidents (
    id          BIGSERIAL PRIMARY KEY,
    monitor_id  BIGINT NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    started_at  TIMESTAMPTZ NOT NULL,
    ended_at    TIMESTAMPTZ,
    last_error  TEXT
);

CREATE TABLE notification_channels (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    kind        TEXT NOT NULL,            -- 'telegram' | 'webhook'
    config      JSONB NOT NULL,           -- chat_id/token/url
    is_active   BOOLEAN NOT NULL DEFAULT TRUE
);
```

Подумай самостоятельно: нужны ли партиции по `checks` (при 1000 мониторов × 1 проверка/минуту — это
≈43 млн записей/месяц). Это хорошее упражнение — **делать не надо**, но понять, где оно сломается, надо.

---

## 6. REST API — базовый контракт

Все ответы — JSON. Ошибки — единый формат `{ "error": { "code": "...", "message": "..." } }`.

```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout

GET    /api/v1/monitors                 # список
POST   /api/v1/monitors                 # создать
GET    /api/v1/monitors/{id}
PATCH  /api/v1/monitors/{id}
DELETE /api/v1/monitors/{id}

GET    /api/v1/monitors/{id}/checks
GET    /api/v1/monitors/{id}/stats?period=24h
GET    /api/v1/monitors/{id}/incidents

GET    /api/v1/channels
POST   /api/v1/channels
DELETE /api/v1/channels/{id}

GET    /healthz
GET    /readyz
GET    /metrics
```

Middleware-стек (в таком порядке):

1. Recover (паники → 500, без падения процесса)
2. Request ID (генерим UUID, кладём в контекст и ответный заголовок)
3. Structured logger (slog + request_id, method, path, status, duration)
4. CORS (если планируется отдельный фронт)
5. Auth (проверка JWT, кладём `user_id` в контекст) — только для приватных роутов
6. Rate limit на `/auth/login` и `/auth/register` (Redis, token bucket per IP)

---

## 7. Docker и docker-compose

`Dockerfile` — обязательно multi-stage:

1. Стадия сборки: `golang:1.22-alpine`, `CGO_ENABLED=0`, статическая сборка.
2. Финальная стадия: `gcr.io/distroless/static` или `alpine:3.19`. Ничего лишнего.
3. Указать `USER nonroot`, `EXPOSE 8080`.

`docker-compose.yml` должен поднимать единой командой `docker compose up`:

- `pulse` — само приложение
- `postgres:16-alpine` с volume, healthcheck
- `redis:7-alpine` с healthcheck
- `migrate` — одноразовый контейнер, прогоняющий миграции перед стартом `pulse`
  (через `depends_on: condition: service_healthy` + `restart: "no"`)
- Опционально: `prometheus` + `grafana` — чтобы увидеть свои метрики в дашборде

Все секреты (DB password, JWT secret, Telegram token) — через `.env`-файл, не хардкод.

---

## 8. Поэтапный план работы

Разбей на этапы и коммить по ходу — так проект доживёт до конца.

### Этап 1. Каркас (1–2 вечера)

- Инициализировать модуль, настроить `Makefile`.
- Поднять docker-compose с Postgres и Redis.
- `main.go` + `/healthz`, slog-логи, чтение конфига из env, graceful shutdown.

### Этап 2. Аутентификация (2–3 вечера)

- Миграции `users` и `refresh_tokens`.
- Регистрация, логин, JWT middleware, refresh/logout.
- Юнит-тесты на сервис аутентификации с моками репозитория.

### Этап 3. CRUD мониторов (1–2 вечера)

- Миграции `monitors`.
- Все эндпоинты, валидация входа, проверка `user_id`.
- Интеграционные тесты на handlers через `httptest` + testcontainers Postgres.

### Этап 4. Воркер и проверки (это самое интересное, 3–5 вечеров)

- Таблица `checks`.
- Планировщик + пул воркеров + HTTP-клиент с таймаутами.
- Graceful shutdown: scheduler, воркеры и HTTP-сервер завершаются через `errgroup` и общий `context`.
- Метрики Prometheus.

### Этап 5. Инциденты и уведомления (2–3 вечера)

- Таблица `incidents`, логика открытия/закрытия.
- Интерфейс `Notifier`, реализация Telegram (или webhook — что проще достать).
- Правило «подтверждение через N подряд fail».

### Этап 6. Агрегаты и статистика (1–2 вечера)

- `GET /monitors/{id}/stats` — uptime %, avg/p95 latency за период.
- Кэширование ответа этой ручки в Redis на 30–60 сек.

### Этап 7. Полировка (1–2 вечера)

- `README.md` с инструкцией запуска и примерами `curl`.
- `golangci-lint` и починка всего, что он скажет.
- GitHub Actions: lint + test + docker build.
- (Опционально) `openapi.yaml` и его рендер в Swagger UI.

Итого ≈ 2–3 недели по вечерам. Если идёт дольше — нормально, ты учишься, а не сдаёшь курсовую.

---

## 9. Что осознанно **не делаем** в MVP

Чтобы проект действительно закончить, зафиксируй «не в скоупе»:

- Никакого фронтенда. API + Postman/`curl` — достаточно. Фронт — отдельный пет-проект при желании.
- Никаких TCP/ICMP/TLS-cert проверок — только HTTP(S).
- Никаких статус-страниц для публичного показа.
- Никаких команд/групп/ролей. Один пользователь = свои мониторы.
- Никакого шардирования/кластера воркеров. Один процесс.
- Никакой очереди вроде Kafka/NATS. Канала + пула горутин в процессе хватит.

Если руки чешутся — занеси в `TODO.md` и продолжай по плану. Это важный навык: отличать «надо» от
«хочется».

---

## 10. Что ты реально отработаешь на этом проекте

Чек-лист навыков, который можно смело положить в резюме после завершения:

- Сборка Go-приложения: layout, `internal`, dependency injection руками.
- Работа с `context.Context`, отмена, таймауты, graceful shutdown через `errgroup`.
- Конкурентность: горутины, каналы, семафор, воркер-пул.
- HTTP-сервер на `net/http`/`chi`, middleware, валидация, унифицированные ошибки.
- Аутентификация: JWT, refresh-токены, хэширование паролей, rate limiting.
- PostgreSQL: миграции, индексы, транзакции, `pgx`, продуманные запросы.
- Redis: кэш, rate limit.
- Внешние интеграции с настоящим HTTP-клиентом (и его тюнинг — transport, timeouts, keep-alive).
- Логирование через `log/slog`, структурные логи, request-id.
- Метрики Prometheus.
- Тестирование: юнит + интеграционные через `testcontainers-go`.
- Docker multi-stage, docker-compose, healthchecks, миграции как отдельный контейнер.
- CI: lint + test + build в GitHub Actions.

---

## 11. Как проверить, что «проект сделан»

Считай проект готовым, когда:

1. `docker compose up` поднимает весь стек с нуля, без ручных шагов (миграции применяются сами).
2. Можно зарегистрироваться, залогиниться, создать 3 монитора (валидный URL, невалидный URL, медленный URL).
3. Через пару минут в БД есть записи `checks`, в `/metrics` ненулевые счётчики, в Telegram пришло
   уведомление о падении и о восстановлении.
4. `go test ./...` проходит зелёным, в том числе интеграционные тесты на testcontainers.
5. `golangci-lint run` молчит.
6. `README.md` написан так, что незнакомый человек (или ты через полгода) запустит проект с нуля.

Удачи. Главный совет: не пытайся сделать красиво сразу. Сделай работающую вертикаль (один монитор
→ одна проверка → одна запись в БД → одно уведомление в лог), а потом расширяй вширь и вглубь.
