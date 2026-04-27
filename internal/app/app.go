package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/RabotaemActivno/pulse/internal/config"
	"github.com/RabotaemActivno/pulse/internal/controller/http"
	"github.com/RabotaemActivno/pulse/pkg/httpserver"
	"github.com/go-chi/chi/v5"
)

func Run(ctx context.Context, cfg config.Config) error {

	r := chi.NewRouter()
	http.PulseRouter(r)
	httpserver.New(r, cfg.HTTP)

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	<-sig

	return nil
}
