package v1

import (
	"net/http"

	"github.com/RabotaemActivno/pulse/pkg/render"
)

func (h *Handlers) GetMonitors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	output, err := h.usecase.GetMonitors(ctx)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "request failed")
		return
	}

	render.JSON(w, output, http.StatusOK)
}
