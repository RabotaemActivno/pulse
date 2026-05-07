package v1

import (
	"encoding/json"
	"net/http"

	"github.com/RabotaemActivno/pulse/internal/dto"
	"github.com/RabotaemActivno/pulse/pkg/render"
)

func (h *Handlers) LogoutUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	input := dto.LogoutUserInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "request failed")
		return
	}

	err = h.usecase.LogoutUser(ctx, input)
	if err != nil {
		render.Error(w, err, http.StatusBadRequest, "request failed")
		return
	}
	
	render.JSON(w, "", http.StatusOK)
}
