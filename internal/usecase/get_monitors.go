package usecase

import (
	"context"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/dto"
)

func (uc *UseCase) GetMonitors(ctx context.Context) (dto.GetMonitorsOutput, error) {

	ms, err := uc.postgres.GetMonitors(ctx)
	if err != nil {
		return dto.GetMonitorsOutput{}, fmt.Errorf("failed to load monitors")
	}

	return dto.GetMonitorsOutput{Monitors: ms}, nil
}
