package usecase

import (
	"context"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/RabotaemActivno/pulse/internal/dto"
	"github.com/RabotaemActivno/pulse/pkg/auth"
	"github.com/rs/zerolog/log"
)

func (uc *UseCase) LoginUser(ctx context.Context, input dto.LoginUserInput) (dto.LoginUserOutput, error) {

	userID, err := uc.postgres.LoginUser(ctx, input.Email, input.Password)
	if err != nil {
		return dto.LoginUserOutput{}, fmt.Errorf("LoginUser: %w", err)
	}

	log.Info().Msg("User logged successful")

	accessToken, err := auth.GenerateAccessToken(userID)
	if err != nil {
		return dto.LoginUserOutput{}, fmt.Errorf("GenerateAccessToken: %w", err)
	}

	refreshToken, err := domain.GenerateRefreshToken()
	if err != nil {
		return dto.LoginUserOutput{}, fmt.Errorf("GenerateRefreshToken: %w", err)
	}

	token, err := domain.NewToken(userID, refreshToken)
	if err != nil {
		return dto.LoginUserOutput{}, fmt.Errorf("NewToken: %w", err)
	}

	err = uc.postgres.SaveRefreshToken(ctx, token.ID, token.UserID, token.TokenHash, token.ExpiresAt)
	if err != nil {
		return dto.LoginUserOutput{}, fmt.Errorf("SaveRefreshToken: %w", err)
	}

	return dto.LoginUserOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
