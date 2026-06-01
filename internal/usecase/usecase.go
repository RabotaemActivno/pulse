package usecase

import (
	"context"
	"time"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/google/uuid"
)

type Postgres interface {
	RegisterUser(ctx context.Context, user domain.User) error
	LoginUser(ctx context.Context, email, password string) (uuid.UUID, error)
	SaveRefreshToken(
		ctx context.Context,
		id uuid.UUID,
		userID uuid.UUID,
		tokenHash string,
		expiresAt time.Time,
	) error
	FindToken(ctx context.Context, tokenHash string) (domain.Token, error)
	LogoutUser(ctx context.Context, tokenHash string) error
	GetMonitors(ctx context.Context) ([]domain.Monitor, error)
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}
