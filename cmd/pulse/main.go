package main

import (
	"context"

	"github.com/RabotaemActivno/pulse/internal/app"
	"github.com/RabotaemActivno/pulse/internal/config"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()

	err := app.Run(ctx, cfg)
	if err != nil {
		panic("error to start App")
	}
}
