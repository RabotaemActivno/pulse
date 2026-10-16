package app

import (
	"context"
	"fmt"
	"log/slog"
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
)

const (
	EnvProd = "prod"
	EnvDev  = "dev"
)

func Run(ctx context.Context, cfg config.Config) error {

	pgPool, err := postgresPool.New(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("postgres New: %w", err)
	}

	log := setupLogger(cfg.Logger)

	p := postgres.New(pgPool)

	uc := usecase.New(p)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	http.PulseRouter(r, uc)
	httpServer := httpserver.New(r, cfg.HTTP)

	sig := make(chan os.Signal, 1)

	log.Info("App started")

	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	<-sig

	httpServer.Close()

	log.Info("App stopped")

	return nil
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case EnvProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}),
		)
	case EnvDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}
