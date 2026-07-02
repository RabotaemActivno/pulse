package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Monitor struct {
	ID                  uuid.UUID `json:"id"`
	UserID              uuid.UUID `json:"user_id"`
	Name                string    `json:"name" validate:"required,min=3,max=64"`
	URL                 string    `json:"URL" validate:"required,url,min=3,max=64"`
	Method              string    `json:"method"`
	IntervalSeconds     int       `json:"interval_seconds" validate:"required,min=1"`
	TimeoutSeconds      int       `json:"timeout_seconds" validate:"required,min=1"`
	ExpectedStatus      int       `json:"expected_status"`
	IsActive            bool      `json:"is_active"`
	NextCheckAt         time.Time `json:"next_check_at"`
	ConsecutiveFailures int       `json:"consecutive_failures"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func NewMonitor(
	UserID uuid.UUID,
	name string,
	url string,
	method string,
	intervalSeconds int,
	timeoutSeconds int,
	expectedStatus int,
) (Monitor, error) {

	if method == "" {
		method = "GET"
	}

	if expectedStatus == 0 {
		expectedStatus = 200
	}

	mtr := Monitor{
		ID:                  uuid.New(),
		UserID:              UserID,
		Name:                name,
		URL:                 url,
		Method:              method,
		IntervalSeconds:     intervalSeconds,
		TimeoutSeconds:      timeoutSeconds,
		ExpectedStatus:      expectedStatus,
		IsActive:            true,
		ConsecutiveFailures: 0,
	}

	return mtr, nil
}

func (m Monitor) Validate() error {
	err := validate.Struct(m)
	if err != nil {
		return fmt.Errorf("validate structL %w", err)
	}

	return nil
}
