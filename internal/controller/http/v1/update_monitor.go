package v1

import (
	"encoding/json"
	"net/http"

	"github.com/RabotaemActivno/pulse/internal/dto"
	"github.com/RabotaemActivno/pulse/pkg/render"
)

func (h *Handlers) UpdateMonitor(w http.ResponseWriter, r *http.Request) {
	input := dto.UpdateMonitorInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "json decode error")
		return
	}

	err = h.usecase.UpdateMonitor(r.Context(), input)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "request field")
		return
	}

	render.Success(w, http.StatusOK, "ok")
}
