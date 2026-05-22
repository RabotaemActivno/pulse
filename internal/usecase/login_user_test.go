package usecase_test

import (
	"context"
	"testing"

	"github.com/RabotaemActivno/pulse/internal/dto"
	"github.com/RabotaemActivno/pulse/internal/usecase"
	"github.com/RabotaemActivno/pulse/internal/usecase/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_Login_Success(t *testing.T) {
	userID := uuid.New()
	postgres := mocks.NewMockPostgres(t)
	postgres.EXPECT().
		LoginUser(mock.Anything, "mock@gmail.com", "123").
		Return(userID, nil)
	postgres.EXPECT().
		SaveRefreshToken(mock.Anything, mock.Anything, userID, mock.Anything, mock.Anything).
		Return(nil)

	uc := usecase.New(postgres)

	input := dto.LoginUserInput{Email: "mock@gmail.com", Password: "123"}

	output, err := uc.LoginUser(context.Background(), input)

	require.NoError(t, err)
	require.NotEmpty(t, output)
}

func Test_Login_EmptyEmail(t *testing.T) {
	postgres := mocks.NewMockPostgres(t)

	uc := usecase.New(postgres)

	input := dto.LoginUserInput{Email: "", Password: "123"}

	output, err := uc.LoginUser(context.Background(), input)

	require.Error(t, err)
	require.Empty(t, output)
}
