package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/RabotaemActivno/pulse/internal/dto"
)

func (uc *UseCase) CreateMonitor(ctx context.Context, input dto.CreateMonitorInput) (dto.CreateMonitorOutput, error) {

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
		return output, fmt.Errorf("domain.CreateMonitor: %w", err)
	}

	err = uc.postgres.CreateMonitor(ctx, mtr)
	if err != nil {
		if errors.Is(err, domain.ErrorMonitorExists) {
			return output, fmt.Errorf("Monitor alredy exists: %w", err)
		}
		return output, fmt.Errorf("db.CreateMonitor: %w", err)
	}

	output.ID = mtr.ID

	return output, nil
}
