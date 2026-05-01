package usecase

import (
	"context"
	"fmt"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/RabotaemActivno/pulse/internal/dto"
)

func (uc *UseCase) RegisterUser(ctx context.Context, userInput dto.RegisterUserInput) (dto.RegisterUserOutput, error) {

	var output dto.RegisterUserOutput

	user, err := domain.NewUser(userInput.Email, userInput.Password)
	if err != nil {
		return output, fmt.Errorf("domain.NewUser: %w", err)
	}

	err = uc.postgres.RegisterUser(ctx, user)
	if err != nil {
		return dto.RegisterUserOutput{}, fmt.Errorf("postgres.RegisterUser: %w", err)
	}

	output.ID = user.ID

	return output, nil
}
