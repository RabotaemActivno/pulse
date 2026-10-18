package v1

import (
	"log/slog"

	"github.com/RabotaemActivno/pulse/internal/usecase"
)

type Handlers struct {
	usecase *usecase.UseCase
	logger  *slog.Logger
}

func New(uc *usecase.UseCase, logger *slog.Logger) *Handlers {
	return &Handlers{usecase: uc, logger: logger}
}
