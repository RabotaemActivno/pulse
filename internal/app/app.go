package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/RabotaemActivno/pulse/internal/adapter/postgres"
	"github.com/RabotaemActivno/pulse/internal/config"
	"github.com/RabotaemActivno/pulse/internal/controller/http"
	"github.com/RabotaemActivno/pulse/pkg/httpserver"
	postgresPool "github.com/RabotaemActivno/pulse/pkg/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func Run(ctx context.Context, cfg config.Config) error {

	pgPool, err := postgresPool.New(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("postgres New: %w", err)
	}
	_ = postgres.New(pgPool)

	r := chi.NewRouter()
	http.PulseRouter(r)
	httpServer := httpserver.New(r, cfg.HTTP)

	log.Info().Msg("App started")
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	<-sig
	log.Info().Msg("App got signal to stop")

	httpServer.Close()

	log.Info().Msg("App stopped")
	return nil
}
