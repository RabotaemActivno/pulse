package postgres

import (
	"context"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/google/uuid"
)

func (p *Postgres) DeleteMonitor(ctx context.Context, monitorID uuid.UUID) error {
	sql := `DELETE FROM monitors WHERE id = $1`
	res, err := p.PgPool.Exec(ctx, sql, monitorID)
	if err != nil {
		return fmt.Errorf("PgPool.Exec: %w", err)
	}
	if res.RowsAffected() == 0 {
		return domain.ErrorNotFound
	}
	return nil
}
