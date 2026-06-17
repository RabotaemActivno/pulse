package usecase

import (
	"context"
	"net/url"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/RabotaemActivno/pulse/internal/dto"
)

func (uc *UseCase) CreateMonitor(ctx context.Context, input dto.CreateMonitorInput) (dto.CreateMonitorOutput, error) {

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
		
	}

}
