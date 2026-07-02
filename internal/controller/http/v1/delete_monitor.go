package v1

import (
	"net/http"

	"github.com/RabotaemActivno/pulse/pkg/render"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type output struct {
	ID uuid.UUID `json:"id"`
}

func (h *Handlers) DeleteMonitor(w http.ResponseWriter, r *http.Request) {

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "invalid monitor id")
		return
	}

	err = h.usecase.DeleteMonitor(r.Context(), id)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "request failed")
		return
	}

	render.JSON(w, output{ID: id}, http.StatusOK)
}
