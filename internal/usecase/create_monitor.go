package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/RabotaemActivno/pulse/internal/dto"
)

func (uc *UseCase) CreateMonitor(ctx context.Context, input dto.CreateMonitorInput) (dto.CreateMonitorOutput, error) {
	const op = "usecase.CreateMonitor"
	log := uc.logger.With("op", op)
	var output dto.CreateMonitorOutput

	mtr, err := domain.NewMonitor(
		input.UserID,
		input.Name,
		input.URL,
		input.Method,
		input.TimeoutSeconds,
		input.IntervalSeconds,
		input.ExpectedStatus,
	)
	if err != nil {
		return output, fmt.Errorf("%v: %w", op, err)
	}

	err = mtr.Validate()
	if err != nil {
		return output, fmt.Errorf("%v: %w", op, err)
	}

	err = uc.postgres.CreateMonitor(ctx, mtr)
	if err != nil {
		if errors.Is(err, domain.ErrorMonitorExists) {
			return output, fmt.Errorf("%v: %w", op, err)
		}
		return output, fmt.Errorf("%v: %w", op, err)
	}

	output.ID = mtr.ID

	log.Info("monitor was created")

	return output, nil
}
