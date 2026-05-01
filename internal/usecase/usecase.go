package usecase

import (
	"context"

	"github.com/RabotaemActivno/pulse/internal/domain"
)

type Postgres interface {
	RegisterUser(ctx context.Context, user domain.User) error
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}
