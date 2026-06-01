package postgres

import (
	"context"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) GetMonitors(ctx context.Context) ([]domain.Monitor, error) {
	sql := `SELECT * FROM monitors`

	rows, err := p.PgPool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("GetMonitors: %w", err)
	}

	monitors, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Monitor])
	if err != nil {
		return nil, fmt.Errorf("GetMonitors: %w", err)
	}

	return monitors, nil
}
