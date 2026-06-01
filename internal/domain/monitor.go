package domain

import (
	"time"

	"github.com/google/uuid"
)

type Monitor struct {
	ID                  uuid.UUID `json:"id"`
	UserID              uuid.UUID `json:"user_id"`
	Name                string    `json:"name"`
	URL                 string    `json:"URL"`
	Method              string    `json:"method"`
	IntervalSeconds     int       `json:"interval_seconds"`
	TimeoutSeconds      int       `json:"timeout_seconds"`
	ExpectedStatus      int       `json:"expected_status"`
	IsActive            bool      `json:"is_active"`
	NextCheckAt         time.Time `json:"next_check_at"`
	ConsecutiveFailures int       `json:"consecutive_failures"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
