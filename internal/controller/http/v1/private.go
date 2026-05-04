package v1

import (
	"net/http"

	"github.com/RabotaemActivno/pulse/pkg/render"
)

func (h *Handlers) Private(w http.ResponseWriter, r *http.Request) {

	render.JSON(w, "private output", http.StatusOK)
}
