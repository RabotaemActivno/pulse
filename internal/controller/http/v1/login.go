package v1

import (
	"encoding/json"
	"net/http"

	"github.com/RabotaemActivno/pulse/internal/dto"
	"github.com/RabotaemActivno/pulse/pkg/render"
)

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	input := dto.LoginUserInput{}
	
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "json decode error")
		return
	}

	output, err := h.usecase.LoginUser(ctx, input)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "request failed")
		return
	}

	render.JSON(w, output, http.StatusOK)
}
