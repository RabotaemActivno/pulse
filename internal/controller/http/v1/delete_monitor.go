package v1

import (
	"errors"
	"net/http"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/RabotaemActivno/pulse/internal/middleware"
	"github.com/RabotaemActivno/pulse/pkg/render"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type output struct {
	ID uuid.UUID `json:"id"`
}

func (h *Handlers) DeleteMonitor(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "invalid monitor id")
		return
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok || userID == uuid.Nil {
		render.Error(w, nil, http.StatusUnauthorized, "unauthorized")
		return
	}

	err = h.usecase.DeleteMonitor(ctx, id, userID)
	if err != nil {
		if errors.Is(err, domain.ErrorNotFound) {
			render.Error(w, err, http.StatusNotFound, "request failed")
		}
		render.Error(w, err, http.StatusBadRequest, "request failed")
		return
	}

	render.JSON(w, output{ID: id}, http.StatusOK)
}
