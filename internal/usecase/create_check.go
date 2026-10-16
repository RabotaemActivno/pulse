package usecase

import (
	"context"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/RabotaemActivno/pulse/pkg/checker"
)

func (uc *UseCase) CreateCheck(ctx context.Context, m domain.Monitor, res checker.CheckResult) error {
	const op = "usecase.CreateCheck"
	check := domain.NewCheck(m, res)
	err := uc.postgres.CreateCheck(ctx, check)
	if err != nil {
		return fmt.Errorf("%v: %w", op, err)
	}

	return nil
}
