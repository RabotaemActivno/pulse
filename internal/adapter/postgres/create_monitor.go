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
		INSERT INTO monitors (id)
	`

	return nil
}
