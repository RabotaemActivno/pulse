package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func PulseRouter(r *chi.Mux) {

	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
