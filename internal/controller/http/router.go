package http

import (
	"net/http"

	ver1 "github.com/RabotaemActivno/pulse/internal/controller/http/v1"
	"github.com/RabotaemActivno/pulse/internal/middleware"
	"github.com/RabotaemActivno/pulse/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func PulseRouter(r *chi.Mux, uc *usecase.UseCase) {

	v1 := ver1.New(uc)

	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		log.Info().Msg("Health checked")
	})

	r.Route("/api", func(r chi.Router) {

		r.Route("/v1", func(r chi.Router) {
			r.Post("/register", v1.CreateUser)
			r.Post("/login", v1.Login)
			r.Post("/refresh", v1.Refresh)
			r.Post("/logout", v1.LogoutUser)

			r.Route("/user", func(r chi.Router) {
				r.Use(middleware.AuthMiddleware)
				r.Get("/private", v1.Private)
			})
		})
	})
}
