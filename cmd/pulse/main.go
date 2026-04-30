package main

import (
	"context"

	"github.com/RabotaemActivno/pulse/internal/app"
	"github.com/RabotaemActivno/pulse/internal/config"
	"github.com/RabotaemActivno/pulse/pkg/logger"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()
	logger.Init(cfg.Logger)

	err := app.Run(ctx, cfg)
	if err != nil {
		panic("error to start App")
	}
}
