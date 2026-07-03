package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/RabotaemActivno/pulse/internal/dto"
)

func (uc *UseCase) UpdateMonitor(ctx context.Context, input dto.UpdateMonitorInput) error {
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("input.Validate: %w", err)
	}

	mtr, err := uc.postgres.GetMonitor(ctx, input.ID)
	if err != nil {
		if errors.Is(err, domain.ErrorNotFound) {
			return fmt.Errorf("Monitor not found: %w", err)
		}
		return fmt.Errorf("uc.GetMonitor: %w", err)
	}

	updatedMonitor := update(mtr, input)
	err = uc.postgres.UpdateMonitor(ctx, updatedMonitor)
	if err != nil {
		if errors.Is(err, domain.ErrorNotFound) {
			return fmt.Errorf("Monitor not found: %w", err)
		}
		return fmt.Errorf("uc.UpdateMonitor: %w", err)
	}
	return nil
}

func update(mtr domain.Monitor, input dto.UpdateMonitorInput) domain.Monitor {
	if input.Name != nil {
		mtr.Name = *input.Name
	}
	if input.URL != nil {
		mtr.URL = *input.URL
	}
	if input.Method != nil {
		mtr.Method = *input.Method
	}
	if input.IntervalSeconds != nil {
		mtr.IntervalSeconds = *input.IntervalSeconds
	}
	if input.TimeoutSeconds != nil {
		mtr.TimeoutSeconds = *input.TimeoutSeconds
	}
	if input.ExpectedStatus != nil {
		mtr.ExpectedStatus = *input.ExpectedStatus
	}
}
