package postgres

import (
	"context"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/dto"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) GetMonitors(ctx context.Context) (dto.GetMonitorsOutput, error) {
	sql := `SELECT * FROM monitors`

	rows, err := p.PgPool.Query(ctx, sql)
	if err != nil {
		return dto.GetMonitorsOutput{}, fmt.Errorf("GetMonitors: %w", err)
	}

	monitors, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.GetMonitorOutput])
	if err != nil {
		return dto.GetMonitorsOutput{}, fmt.Errorf("GetMonitors: %w", err)
	}

	return dto.GetMonitorsOutput{Monitors: monitors}, nil
}
