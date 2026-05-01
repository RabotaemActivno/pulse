package postgres

import (
	"context"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/domain"
)

func (p *Postgres) RegisterUser(ctx context.Context, user domain.User) error {

	queryIsExists := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE email = $1
		)
	`
	var exists bool

	err := p.PgPool.QueryRow(ctx, queryIsExists, user.Email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("Register User: %w", err)
	}

	if exists {
		return fmt.Errorf("User is: %w", domain.ErrorUserExists)
	}

	queryRegisterUser := `
		INSERT INTO users (id, email, password_hash) 
			VALUES ($1, $2, $3)
	`

	args := []any{
		user.ID,
		user.Email,
		user.PasswordHash,
	}

	_, err = p.PgPool.Exec(ctx, queryRegisterUser, args...)
	if err != nil {
		return fmt.Errorf("PgPoo.Exec: %w", err)
	}

	return nil
}
