package usecase

import (
	"context"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/RabotaemActivno/pulse/internal/dto"
)

func (uc *UseCase) LogoutUser(ctx context.Context, input dto.LogoutUserInput) error {

	tokenHash := domain.RefreshTokenToHash(input.RefreshToken)

	err := uc.postgres.LogoutUser(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("db.LogoutUser: %w", err)
	}

	return nil
}
