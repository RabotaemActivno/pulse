package postgres

import (
	"context"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/google/uuid"
)

func (p *Postgres) LoginUser(ctx context.Context, email string, password string) (uuid.UUID, error) {
	sql := `
		SELECT id, email, password_hash, created_at FROM users WHERE email = $1
	`

	var usr domain.User

	err := p.PgPool.QueryRow(ctx, sql, email).Scan(&usr.ID, &usr.Email, &usr.PasswordHash, &usr.CreatedAt)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Login user: %w", err)
	}

	if err := usr.CheckPassword(password); err != nil {
		return uuid.UUID{}, fmt.Errorf(domain.ErrorCredentials.Error(), err)
	}

	return usr.ID, nil
}
