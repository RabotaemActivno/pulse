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
	"github.com/RabotaemActivno/pulse/internal/usecase"
	"github.com/RabotaemActivno/pulse/pkg/httpserver"
	postgresPool "github.com/RabotaemActivno/pulse/pkg/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func Run(ctx context.Context, cfg config.Config) error {

	pgPool, err := postgresPool.New(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("postgres New: %w", err)
	}

	p := postgres.New(pgPool)

	uc := usecase.New(p)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	http.PulseRouter(r, uc)
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
