package usecase_test

import (
	"context"
	"testing"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/RabotaemActivno/pulse/internal/dto"
	"github.com/RabotaemActivno/pulse/internal/usecase"
	"github.com/RabotaemActivno/pulse/internal/usecase/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_CreateMonitor_Success(t *testing.T) {
	postgres := mocks.NewMockPostgres(t)
	ctx := context.Background()
	userID := uuid.New()
	ctx = context.WithValue(ctx, "user_id", userID)
	mtr, _ := domain.NewMonitor(userID, "Name", "http://mocksite.com", "GET", 1, 1, 200)
	postgres.EXPECT().CreateMonitor(ctx, mock.Anything).Return(nil)

	uc := usecase.New(postgres)

	input := dto.CreateMonitorInput{UserID: mtr.UserID, Name: mtr.Name, URL: mtr.URL, Method: mtr.Method, IntervalSeconds: mtr.IntervalSeconds, TimeoutSeconds: mtr.TimeoutSeconds, ExpectedStatus: mtr.ExpectedStatus}

	output, err := uc.CreateMonitor(ctx, input)

	require.NoError(t, err)
	require.NotEmpty(t, output)
}
