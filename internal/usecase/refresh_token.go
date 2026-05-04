package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/RabotaemActivno/pulse/internal/dto"
	"github.com/RabotaemActivno/pulse/pkg/auth"
)

func (uc *UseCase) RefreshToken(ctx context.Context, input dto.RefreshTokenInput) (dto.RefreshTokenOutput, error) {
	hashedIncomingTkn := domain.RefreshTokenToHash(input.Token)

	tkn, err := uc.postgres.FindToken(ctx, hashedIncomingTkn)
	if err != nil {
		return dto.RefreshTokenOutput{}, fmt.Errorf("db.FindToken: %w", err)
	}

	if tkn.RevokedAt != nil && !tkn.RevokedAt.Before(time.Now()) && tkn.HasExpired() {
		return dto.RefreshTokenOutput{}, fmt.Errorf("Invalid token")
	}

	newTknStr, err := domain.GenerateRefreshToken()
	if err != nil {
		return dto.RefreshTokenOutput{}, fmt.Errorf("domain.GenerateRefreshToken: %w", err)
	}

	newTkn, err := domain.NewToken(tkn.UserID, newTknStr)
	if err != nil {
		return dto.RefreshTokenOutput{}, fmt.Errorf("domain.NewToken: %w", err)
	}

	err = uc.postgres.SaveRefreshToken(ctx, newTkn.ID, newTkn.UserID, newTkn.TokenHash, newTkn.ExpiresAt)
	if err != nil {
		return dto.RefreshTokenOutput{}, fmt.Errorf("uc.SaveRefreshToken: %w", err)
	}

	accessToken, err := auth.GenerateAccessToken(newTkn.UserID)
	if err != nil {
		return dto.RefreshTokenOutput{}, fmt.Errorf("GenerateAccessToken: %w", err)
	}

	return dto.RefreshTokenOutput{
		RefreshToken: newTknStr,
		AccessToken:  accessToken,
	}, nil
}
