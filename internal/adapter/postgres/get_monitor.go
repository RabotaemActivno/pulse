package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) GetMonitor(ctx context.Context, id uuid.UUID) (domain.Monitor, error) {
	sql := `SELECT user_id, 
			name, 
			url, 
			method, 
			interval_seconds, 
			timeout_seconds, 
			expected_status, 
			is_active FROM monitors WHERE id = $1`

	var mtr domain.Monitor
	mtr.ID = id

	err := p.PgPool.QueryRow(ctx, sql, id).Scan(
		&mtr.UserID,
		&mtr.Name,
		&mtr.URL,
		&mtr.Method,
		&mtr.IntervalSeconds,
		&mtr.TimeoutSeconds,
		&mtr.ExpectedStatus,
		&mtr.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Monitor{}, domain.ErrorNotFound
		}
		return domain.Monitor{}, fmt.Errorf("PgPool.Scan: %w", err)
	}

	return mtr, nil
}
