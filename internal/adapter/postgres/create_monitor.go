package postgres

import (
	"context"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/domain"
)

func (p *Postgres) CreateMonitor(ctx context.Context, mtr domain.Monitor) error {

	queryIsExists := `
		SELECT EXISTS (
			SELECT 1
			FROM monitors
			WHERE url = $1
		)
	`
	var exists bool
	err := p.PgPool.QueryRow(ctx, queryIsExists, mtr.URL).Scan(&exists)
	if err != nil {
		return fmt.Errorf("Create Monitor: %w", err)
	}
	if exists {
		return domain.ErrorMonitorExists
	}

	queryCreateMonitor := `
		INSERT INTO monitors (
			id, 
			user_id, 
			name, 
			url, 
			method, 
			interval_seconds, 
			timeout_seconds,
			expected_status,
		  	consecutive_failures
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	args := []any{
		mtr.ID,
		mtr.UserID,
		mtr.Name,
		mtr.URL,
		mtr.Method,
		mtr.IntervalSeconds,
		mtr.TimeoutSeconds,
		mtr.ExpectedStatus,
		mtr.ConsecutiveFailures,
	}

	_, err = p.PgPool.Exec(ctx, queryCreateMonitor, args...)
	if err != nil {
		return fmt.Errorf("PgPool.Exec: %w", err)
	}

	return nil
}
