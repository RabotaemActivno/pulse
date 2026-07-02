package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) UpdateMonitor(ctx context.Context, mtr domain.Monitor) error {
	sql := `UPDATE monitors SET 
                    name = $1, 
                    url = $2, 
                    method = $3, 
                    interval_seconds = $4, 
                    timeout_seconds = $5, 
                    expected_status = $6,
                    updated_at = NOW() WHERE id = $7`

	args := []any{
		mtr.Name,
		mtr.URL,
		mtr.Method,
		mtr.IntervalSeconds,
		mtr.TimeoutSeconds,
		mtr.ExpectedStatus,
		mtr.ID,
	}

	_, err := p.PgPool.Exec(ctx, sql, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrorNotFound
		}
		return fmt.Errorf("pgPool, UpdateMonitor: %w", err)
	}

	return nil
}
