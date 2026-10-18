package postgres

import (
	"context"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/google/uuid"
)

func (p *Postgres) DeleteMonitor(ctx context.Context, monitorID uuid.UUID, userID uuid.UUID) error {
	sql := `DELETE FROM monitors WHERE id = $1 AND user_id = $2`
	res, err := p.PgPool.Exec(ctx, sql, monitorID, userID)
	if err != nil {
		return fmt.Errorf("PgPool.Exec: %w", err)
	}
	if res.RowsAffected() == 0 {
		return domain.ErrorNotFound
	}
	return nil
}
