package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/RabotaemActivno/pulse/internal/dto"
	"github.com/RabotaemActivno/pulse/internal/usecase"
	"github.com/RabotaemActivno/pulse/internal/usecase/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_Refresh_Success(t *testing.T) {
	input := dto.RefreshTokenInput{Token: "123"}
	userID := uuid.New()
	hashedToken := domain.RefreshTokenToHash(input.Token)
	tokenFromPgs, err := domain.NewToken(userID, input.Token)
	postgres := mocks.NewMockPostgres(t)

	postgres.EXPECT().
		FindToken(mock.Anything, hashedToken).
		Return(tokenFromPgs, nil)
	postgres.EXPECT().
		SaveRefreshToken(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	uc := usecase.New(postgres)

	output, err := uc.RefreshToken(context.Background(), input)

	require.NoError(t, err)
	require.NotEmpty(t, output)
}

func Test_Refresh_Expired(t *testing.T) {
	input := dto.RefreshTokenInput{Token: "123"}
	hashedToken := domain.RefreshTokenToHash(input.Token)
	tokenFromPgs := domain.Token{TokenHash: hashedToken, ExpiresAt: time.Now()}
	postgres := mocks.NewMockPostgres(t)

	postgres.EXPECT().
		FindToken(mock.Anything, hashedToken).
		Return(tokenFromPgs, nil)

	uc := usecase.New(postgres)

	output, err := uc.RefreshToken(context.Background(), input)

	require.Error(t, err)
	require.Empty(t, output)
}
