package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/RabotaemActivno/pulse/internal/dto"
	"github.com/RabotaemActivno/pulse/internal/middleware"
	"github.com/RabotaemActivno/pulse/pkg/render"
)

func (h *Handlers) CreateMonitor(w http.ResponseWriter, r *http.Request) {
	input := dto.CreateMonitorInput{}

	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "json decode error")
		return
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		render.Error(w, errors.New("user_id is invalid"), http.StatusBadRequest, "request failed")
		return
	}

	input.UserID = userID

	output, err := h.usecase.CreateMonitor(ctx, input)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "request failed")
		return
	}

	render.JSON(w, output, http.StatusOK)
}
