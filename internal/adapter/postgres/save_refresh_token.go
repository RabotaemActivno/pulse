package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (p *Postgres) SaveRefreshToken(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
	tokenHash string,
	expiresAt time.Time,
) error {
	sql := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE SET
			id = EXCLUDED.id,
			token_hash = EXCLUDED.token_hash,
			expires_at = EXCLUDED.expires_at
	`
	args := []any{
		id,
		userID,
		tokenHash,
		expiresAt,
	}

	_, err := p.PgPool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PgPool: %w", err)
	}

	return nil
}
