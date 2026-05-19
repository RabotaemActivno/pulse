# Pulse — скелет расширения проекта

Документ описывает, как нарастить текущую (слойную) раскладку под мониторы,
проверки, инциденты, планировщик и уведомления. Не повторяет md-описание
продукта (см. `pulse-uptime-monitor.md`), фокус — структура кода.

## 1. Файловое дерево после расширения

```
internal/
├── app/
│   └── app.go                       # composition root: pg, redis, uc, scheduler, checker, notifier, http
├── config/
│   └── config.go                    # + блоки Scheduler, Checker, Notifier, Redis
│
├── domain/
│   ├── errors.go
│   ├── user.go                      # уже есть
│   ├── token.go                     # уже есть
│   ├── monitor.go                   # NEW
│   ├── check.go                     # NEW
│   ├── incident.go                  # NEW
│   └── channel.go                   # NEW (notification_channels)
│
├── dto/
│   ├── login_user.go ... (текущие)
│   ├── create_monitor.go            # NEW
│   ├── update_monitor.go            # NEW
│   ├── list_monitors.go             # NEW
│   ├── list_checks.go               # NEW
│   ├── monitor_stats.go             # NEW
│   ├── list_incidents.go            # NEW
│   └── channel.go                   # NEW
│
├── usecase/
│   ├── repo.go                      # NEW: узкие интерфейсы вместо одного Postgres
│   ├── usecase.go                   # переписать: UseCase агрегирует мини-сервисы или хранит все repo
│   ├── register_user.go ...         # уже есть
│   ├── create_monitor.go            # NEW
│   ├── list_monitors.go             # NEW
│   ├── get_monitor.go               # NEW
│   ├── update_monitor.go            # NEW
│   ├── delete_monitor.go            # NEW
│   ├── list_checks.go               # NEW
│   ├── monitor_stats.go             # NEW
│   ├── list_incidents.go            # NEW
│   ├── create_channel.go            # NEW
│   ├── list_channels.go             # NEW
│   └── delete_channel.go            # NEW
│
├── adapter/
│   ├── postgres/
│   │   ├── postgres.go              # уже есть
│   │   ├── register_user.go ...     # уже есть
│   │   ├── monitor_create.go        # NEW
│   │   ├── monitor_list.go          # NEW
│   │   ├── monitor_get.go           # NEW
│   │   ├── monitor_update.go        # NEW
│   │   ├── monitor_delete.go        # NEW
│   │   ├── monitor_due.go           # NEW: SELECT ... WHERE next_check_at <= now()
│   │   ├── monitor_reschedule.go    # NEW: UPDATE next_check_at, consecutive_failures
│   │   ├── check_insert.go          # NEW
│   │   ├── check_list.go            # NEW
│   │   ├── check_stats.go           # NEW: uptime%, p95, avg
│   │   ├── incident_open.go         # NEW
│   │   ├── incident_close.go        # NEW
│   │   ├── incident_list.go         # NEW
│   │   ├── channel_create.go        # NEW
│   │   ├── channel_list.go          # NEW
│   │   └── channel_delete.go        # NEW
│   ├── redis/                       # NEW
│   │   ├── redis.go                 # клиент
│   │   ├── stats_cache.go           # кэш ответа /stats
│   │   └── ratelimit.go             # token-bucket для /auth/*
│   ├── telegram/                    # NEW
│   │   └── telegram.go              # реализация Notifier
│   └── webhook/                     # NEW
│       └── webhook.go               # реализация Notifier
│
├── scheduler/                       # NEW: тикает, выгребает due-мониторы, кидает в канал
│   └── scheduler.go
│
├── checker/                         # NEW: пул воркеров, http.Client, выполняет пробы
│   ├── checker.go
│   ├── client.go                    # настроенный *http.Client (Transport, keep-alive)
│   └── probe.go                     # одна проба: запрос -> domain.Check
│
├── notifier/                        # NEW: интерфейс + диспетчер по каналам
│   ├── notifier.go
│   └── dispatcher.go                # принимает событие, шлёт во все активные каналы юзера
│
├── middleware/
│   └── auth.go                      # уже есть; добавится rate_limit.go, request_id.go, logger.go
│
└── controller/http/
    ├── router.go                    # + /monitors, /channels
    └── v1/
        ├── v1.go                    # уже есть
        ├── create_user.go ...       # уже есть
        ├── create_monitor.go        # NEW
        ├── list_monitors.go         # NEW
        ├── get_monitor.go           # NEW
        ├── update_monitor.go        # NEW
        ├── delete_monitor.go        # NEW
        ├── list_checks.go           # NEW
        ├── monitor_stats.go         # NEW
        ├── list_incidents.go        # NEW
        ├── create_channel.go        # NEW
        ├── list_channels.go         # NEW
        └── delete_channel.go        # NEW

migration/
├── 0001_users.sql                   # уже есть
├── 0002_refresh_tokens.sql          # уже есть
├── 0003_monitors.sql                # NEW
├── 0004_checks.sql                  # NEW
├── 0005_incidents.sql               # NEW
└── 0006_notification_channels.sql   # NEW

pkg/
├── auth/         (есть)
├── httpserver/   (есть)
├── logger/       (есть)
├── postgres/     (есть)
├── render/       (есть)
└── redis/        # NEW: тонкая обёртка над go-redis (по аналогии с pkg/postgres)
```

---

## 2. Domain — скелеты типов

`internal/domain/monitor.go`

```go
package domain

import (
    "time"

    "github.com/google/uuid"
)

type HTTPMethod string

const (
    MethodGET  HTTPMethod = "GET"
    MethodHEAD HTTPMethod = "HEAD"
    MethodPOST HTTPMethod = "POST"
)

type Monitor struct {
    ID                  uuid.UUID  `json:"id"`
    UserID              uuid.UUID  `json:"user_id"`
    Name                string     `json:"name"             validate:"required,min=1,max=120"`
    URL                 string     `json:"url"              validate:"required,url"`
    Method              HTTPMethod `json:"method"           validate:"oneof=GET HEAD POST"`
    IntervalSeconds     int        `json:"interval_seconds" validate:"min=10,max=3600"`
    TimeoutSeconds      int        `json:"timeout_seconds"  validate:"min=1,max=30"`
    ExpectedStatus      int        `json:"expected_status"  validate:"min=100,max=599"`
    IsActive            bool       `json:"is_active"`
    NextCheckAt         time.Time  `json:"next_check_at"`
    ConsecutiveFailures int        `json:"consecutive_failures"`
    CreatedAt           time.Time  `json:"created_at"`
    UpdatedAt           time.Time  `json:"updated_at"`
}

func NewMonitor(userID uuid.UUID, name, url string, method HTTPMethod,
    interval, timeout, expectedStatus int) (Monitor, error) {
    m := Monitor{
        ID:              uuid.New(),
        UserID:          userID,
        Name:            name,
        URL:             url,
        Method:          method,
        IntervalSeconds: interval,
        TimeoutSeconds:  timeout,
        ExpectedStatus:  expectedStatus,
        IsActive:        true,
        NextCheckAt:     time.Now(),
    }
    if err := validate.Struct(m); err != nil {
        return Monitor{}, err
    }
    return m, nil
}
```

`internal/domain/check.go`

```go
package domain

import (
    "time"

    "github.com/google/uuid"
)

type CheckStatus string

const (
    CheckUp   CheckStatus = "up"
    CheckDown CheckStatus = "down"
)

type Check struct {
    ID         uuid.UUID   `json:"id"`
    MonitorID  uuid.UUID   `json:"monitor_id"`
    CheckedAt  time.Time   `json:"checked_at"`
    Status     CheckStatus `json:"status"`
    HTTPStatus *int        `json:"http_status,omitempty"`
    LatencyMs  *int        `json:"latency_ms,omitempty"`
    Error      *string     `json:"error,omitempty"`
}
```

`internal/domain/incident.go`

```go
package domain

import (
    "time"

    "github.com/google/uuid"
)

type Incident struct {
    ID        uuid.UUID  `json:"id"`
    MonitorID uuid.UUID  `json:"monitor_id"`
    StartedAt time.Time  `json:"started_at"`
    EndedAt   *time.Time `json:"ended_at,omitempty"`
    LastError *string    `json:"last_error,omitempty"`
}

func (i Incident) IsOpen() bool { return i.EndedAt == nil }
```

`internal/domain/channel.go`

```go
package domain

import (
    "encoding/json"

    "github.com/google/uuid"
)

type ChannelKind string

const (
    ChannelTelegram ChannelKind = "telegram"
    ChannelWebhook  ChannelKind = "webhook"
)

type Channel struct {
    ID       uuid.UUID       `json:"id"`
    UserID   uuid.UUID       `json:"user_id"`
    Kind     ChannelKind     `json:"kind"`
    Config   json.RawMessage `json:"config"`
    IsActive bool            `json:"is_active"`
}
```

---

## 3. Usecase — разбиение интерфейса репозитория

`internal/usecase/repo.go`

```go
package usecase

import (
    "context"
    "time"

    "github.com/RabotaemActivno/pulse/internal/domain"
    "github.com/google/uuid"
)

type UserRepo interface {
    RegisterUser(ctx context.Context, u domain.User) error
    LoginUser(ctx context.Context, email, password string) (uuid.UUID, error)
}

type TokenRepo interface {
    SaveRefreshToken(ctx context.Context, id, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
    FindToken(ctx context.Context, tokenHash string) (domain.Token, error)
    LogoutUser(ctx context.Context, tokenHash string) error
}

type MonitorRepo interface {
    Create(ctx context.Context, m domain.Monitor) error
    Get(ctx context.Context, userID, id uuid.UUID) (domain.Monitor, error)
    List(ctx context.Context, userID uuid.UUID) ([]domain.Monitor, error)
    Update(ctx context.Context, m domain.Monitor) error
    Delete(ctx context.Context, userID, id uuid.UUID) error

    // для воркера
    ListDue(ctx context.Context, limit int) ([]domain.Monitor, error)
    Reschedule(ctx context.Context, id uuid.UUID, nextAt time.Time, failures int) error
}

type CheckRepo interface {
    Insert(ctx context.Context, c domain.Check) error
    List(ctx context.Context, monitorID uuid.UUID, from, to time.Time, limit int) ([]domain.Check, error)
    Stats(ctx context.Context, monitorID uuid.UUID, period time.Duration) (domain.Stats, error)
}

type IncidentRepo interface {
    OpenIfNone(ctx context.Context, monitorID uuid.UUID, at time.Time, lastErr string) error
    CloseOpen(ctx context.Context, monitorID uuid.UUID, at time.Time) error
    List(ctx context.Context, monitorID uuid.UUID) ([]domain.Incident, error)
}

type ChannelRepo interface {
    Create(ctx context.Context, c domain.Channel) error
    List(ctx context.Context, userID uuid.UUID) ([]domain.Channel, error)
    Delete(ctx context.Context, userID, id uuid.UUID) error
}
```

`internal/usecase/usecase.go` (новый вид)

```go
package usecase

type UseCase struct {
    Users     UserRepo
    Tokens    TokenRepo
    Monitors  MonitorRepo
    Checks    CheckRepo
    Incidents IncidentRepo
    Channels  ChannelRepo
}

type Deps struct {
    Users     UserRepo
    Tokens    TokenRepo
    Monitors  MonitorRepo
    Checks    CheckRepo
    Incidents IncidentRepo
    Channels  ChannelRepo
}

func New(d Deps) *UseCase { return &UseCase{
    Users: d.Users, Tokens: d.Tokens, Monitors: d.Monitors,
    Checks: d.Checks, Incidents: d.Incidents, Channels: d.Channels,
} }
```

`*adapter/postgres.Postgres` реализует все шесть интерфейсов сразу — в `app.go`
один объект передаётся во все поля `Deps`.

Пример usecase для одного действия — `internal/usecase/create_monitor.go`:

```go
package usecase

import (
    "context"

    "github.com/RabotaemActivno/pulse/internal/domain"
    "github.com/RabotaemActivno/pulse/internal/dto"
    "github.com/google/uuid"
)

func (uc *UseCase) CreateMonitor(ctx context.Context, userID uuid.UUID, in dto.CreateMonitor) (domain.Monitor, error) {
    m, err := domain.NewMonitor(userID, in.Name, in.URL, domain.HTTPMethod(in.Method),
        in.IntervalSeconds, in.TimeoutSeconds, in.ExpectedStatus)
    if err != nil {
        return domain.Monitor{}, err
    }
    if err := uc.Monitors.Create(ctx, m); err != nil {
        return domain.Monitor{}, err
    }
    return m, nil
}
```

---

## 4. Scheduler / Checker / Notifier

`internal/scheduler/scheduler.go`

```go
package scheduler

import (
    "context"
    "time"

    "github.com/RabotaemActivno/pulse/internal/domain"
    "github.com/RabotaemActivno/pulse/internal/usecase"
)

type Scheduler struct {
    repo  usecase.MonitorRepo
    out   chan<- domain.Monitor
    tick  time.Duration
    batch int
}

func New(r usecase.MonitorRepo, out chan<- domain.Monitor, tick time.Duration, batch int) *Scheduler {
    return &Scheduler{repo: r, out: out, tick: tick, batch: batch}
}

func (s *Scheduler) Run(ctx context.Context) error {
    t := time.NewTicker(s.tick)
    defer t.Stop()
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-t.C:
            due, err := s.repo.ListDue(ctx, s.batch)
            if err != nil {
                continue // log
            }
            for _, m := range due {
                select {
                case s.out <- m:
                case <-ctx.Done():
                    return ctx.Err()
                }
            }
        }
    }
}
```

`internal/checker/checker.go`

```go
package checker

import (
    "context"
    "sync"

    "github.com/RabotaemActivno/pulse/internal/domain"
    "github.com/RabotaemActivno/pulse/internal/usecase"
)

type Checker struct {
    in        <-chan domain.Monitor
    workers   int
    monitors  usecase.MonitorRepo
    checks    usecase.CheckRepo
    incidents usecase.IncidentRepo
    notifyCh  chan<- domain.Incident
    probe     Prober
}

type Prober interface {
    Do(ctx context.Context, m domain.Monitor) domain.Check
}

func (c *Checker) Run(ctx context.Context) error {
    var wg sync.WaitGroup
    for i := 0; i < c.workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            c.loop(ctx)
        }()
    }
    wg.Wait()
    return nil
}

func (c *Checker) loop(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case m, ok := <-c.in:
            if !ok {
                return
            }
            // probe -> insert check -> reschedule -> open/close incident -> push to notifier
            _ = c.probe.Do(ctx, m)
        }
    }
}
```

`internal/notifier/notifier.go`

```go
package notifier

import (
    "context"

    "github.com/RabotaemActivno/pulse/internal/domain"
)

type Event struct {
    Monitor  domain.Monitor
    Incident domain.Incident
    Kind     EventKind // "down" | "up"
}

type EventKind string

const (
    Down EventKind = "down"
    Up   EventKind = "up"
)

type Sender interface {
    Send(ctx context.Context, ch domain.Channel, e Event) error
}

type Dispatcher struct {
    senders  map[domain.ChannelKind]Sender
    channels ChannelLister
    in       <-chan Event
}

type ChannelLister interface {
    ByUser(ctx context.Context, userID string) ([]domain.Channel, error)
}

func (d *Dispatcher) Run(ctx context.Context) error { /* range d.in -> fanout */ return nil }
```

---

## 5. App: composition root

`internal/app/app.go` (форма после расширения)

```go
func Run(ctx context.Context, cfg config.Config) error {
    pgPool, err := postgresPool.New(ctx, cfg.Postgres)
    if err != nil { return err }

    rdb, err := redis.New(ctx, cfg.Redis)
    if err != nil { return err }

    pg := postgres.New(pgPool) // реализует все *Repo интерфейсы

    uc := usecase.New(usecase.Deps{
        Users: pg, Tokens: pg, Monitors: pg,
        Checks: pg, Incidents: pg, Channels: pg,
    })

    jobs := make(chan domain.Monitor, cfg.Checker.QueueSize)
    events := make(chan notifier.Event, 1024)

    sch := scheduler.New(pg, jobs, cfg.Scheduler.Tick, cfg.Scheduler.Batch)
    chk := checker.New(jobs, events, pg, pg, pg, cfg.Checker)
    ntf := notifier.NewDispatcher(events, pg, telegramSender, webhookSender)

    r := chi.NewRouter()
    http.PulseRouter(r, uc)
    httpServer := httpserver.New(r, cfg.HTTP)

    g, gctx := errgroup.WithContext(ctx)
    g.Go(func() error { return sch.Run(gctx) })
    g.Go(func() error { return chk.Run(gctx) })
    g.Go(func() error { return ntf.Run(gctx) })
    g.Go(func() error { return httpServer.Run() })

    g.Go(func() error {
        sig := make(chan os.Signal, 1)
        signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
        select {
        case <-gctx.Done():
            return gctx.Err()
        case <-sig:
            httpServer.Close()
            return nil
        }
    })

    return g.Wait()
}
```

---

## 6. Маршруты под мониторы (router.go)

```go
r.Route("/api/v1", func(r chi.Router) {
    r.Route("/auth", func(r chi.Router) { /* существующие */ })

    r.Group(func(r chi.Router) {
        r.Use(middleware.AuthMiddleware)

        r.Route("/monitors", func(r chi.Router) {
            r.Get("/",     v1.ListMonitors)
            r.Post("/",    v1.CreateMonitor)
            r.Get("/{id}", v1.GetMonitor)
            r.Patch("/{id}",  v1.UpdateMonitor)
            r.Delete("/{id}", v1.DeleteMonitor)

            r.Get("/{id}/checks",    v1.ListChecks)
            r.Get("/{id}/stats",     v1.MonitorStats)
            r.Get("/{id}/incidents", v1.ListIncidents)
        })

        r.Route("/channels", func(r chi.Router) {
            r.Get("/",        v1.ListChannels)
            r.Post("/",       v1.CreateChannel)
            r.Delete("/{id}", v1.DeleteChannel)
        })
    })
})
```

---

## 7. Порядок реализации

1. Миграции 0003..0006.
2. `domain/monitor.go` + узкие интерфейсы в `usecase/repo.go`, переписать `usecase.go`.
3. CRUD для `monitors`: dto -> usecase -> postgres -> handler -> роуты. Получить рабочий API.
4. `domain/check.go`, `check_*` в postgres, `ListChecks`/`MonitorStats`.
5. `scheduler` + `checker` + интеграция в `app.go` через `errgroup`. Запуск без уведомлений, проверки пишутся в БД.
6. `domain/incident.go`, открытие/закрытие в воркере при N подряд fail.
7. `domain/channel.go`, CRUD каналов, `notifier` + telegram/webhook senders.
8. Redis: кэш `/stats`, rate-limit на `/auth/*`.
9. Метрики Prometheus, `/healthz`, `/readyz`.
