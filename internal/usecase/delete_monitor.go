package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (uc *UseCase) DeleteMonitor(ctx context.Context, monitorID uuid.UUID, userID uuid.UUID) error {
	const op = "usecase.DeleteMonitor"
	log := uc.logger.With("op", op)
	err := uc.postgres.DeleteMonitor(ctx, monitorID, userID)
	if err != nil {
		return fmt.Errorf("%v: %w", op, err)
	}
	log.Info("monitor was deleted")
	return nil
}
