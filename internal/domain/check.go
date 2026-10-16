package domain

import (
	"time"

	"github.com/RabotaemActivno/pulse/pkg/checker"
	"github.com/google/uuid"
)

type CheckStatus string

const (
	CheckStatusUp   CheckStatus = "up"
	CheckStatusDown CheckStatus = "down"
)

//type CheckResult struct {
//	CheckedAt  time.Time
//	HTTPStatus int
//	Latency    time.Duration
//	Err        error
//}

type Check struct {
	ID         uuid.UUID     `json:"id"`
	MonitorID  uuid.UUID     `json:"monitor_id"`
	CheckedAt  time.Time     `json:"checked_at"`
	Status     CheckStatus   `json:"status"`
	HTTPStatus *int          `json:"http_status"`
	Latency    time.Duration `json:"latency"`
	Error      *string       `json:"error"`
}

func NewCheck(m Monitor, r checker.CheckResult) Check {

	if r.CheckedAt.IsZero() {
		r.CheckedAt = time.Now()
	}

	c := Check{
		ID:        uuid.New(),
		MonitorID: m.ID,
		CheckedAt: r.CheckedAt,
		Status:    CheckStatusUp,
		Latency:   r.Latency,
	}

	if r.Err != nil || r.HTTPStatus != m.ExpectedStatus {
		c.Status = CheckStatusDown
	}

	if r.HTTPStatus != 0 {
		c.HTTPStatus = &r.HTTPStatus
	}

	if r.Err != nil {
		c.Error = new(r.Err.Error())
	}
	return c
}
