package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/RabotaemActivno/pulse/internal/dto"
	"github.com/RabotaemActivno/pulse/pkg/render"
	"github.com/google/uuid"
)

func (h *Handlers) CreateMonitor(w http.ResponseWriter, r *http.Request) {
	input := dto.CreateMonitorInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "json decode error")
		return
	}

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		render.Error(w, errors.New("user_id is invalid"), http.StatusBadRequest, "request failed")
		return
	}

	input.UserID = userID

	output, err := h.usecase.CreateMonitor(r.Context(), input)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "request failed")
		return
	}

	render.JSON(w, output, http.StatusOK)
}
