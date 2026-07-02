package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/google/uuid"
)

func (uc *UseCase) DeleteMonitor(ctx context.Context, monitorID uuid.UUID) error {
	err := uc.postgres.DeleteMonitor(ctx, monitorID)
	if err != nil {
		if errors.Is(err, domain.ErrorNotFound) {
			return fmt.Errorf("uc.DeleteMonitor: %w", err)
		}
		return fmt.Errorf("uc.DeleteMonitor")
	}
	return nil
}
