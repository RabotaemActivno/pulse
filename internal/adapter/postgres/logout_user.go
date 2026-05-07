package postgres

import (
	"context"
	"fmt"
	"time"
)

func (p *Postgres) LogoutUser(ctx context.Context, tokenHash string) error {

	sql := `
		UPDATE refresh_tokens SET revoked_at = $1 WHERE token_hash = $2
	`

	_, err := p.PgPool.Exec(ctx, sql, time.Now(), tokenHash)
	if err != nil {
		return fmt.Errorf("db.LogoutUser: %w", err)
	}

	return nil
}
