package v1

import (
	"net/http"

	"github.com/RabotaemActivno/pulse/pkg/render"
)

func (h *Handlers) Private(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()

	render.JSON(w, "", http.StatusOK)
}
