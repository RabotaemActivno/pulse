package postgres

import (
	"context"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/domain"
)

func (p *Postgres) FindToken(ctx context.Context, tokenHash string) (domain.Token, error) {
	sql := `
		SELECT id, user_id, expires_at, revoked_at FROM refresh_tokens WHERE token_hash = $1	
	`

	var tkn domain.Token

	err := p.PgPool.QueryRow(ctx, sql, tokenHash).Scan(&tkn.ID, &tkn.UserID, &tkn.ExpiresAt, &tkn.RevokedAt)
	if err != nil {
		return domain.Token{}, fmt.Errorf("FindToken: %w", err)
	}

	return tkn, nil
}
