package v1

import (
	"net/http"

	"github.com/eulerbutcooler/virtus/internal/handler/response"
	"github.com/eulerbutcooler/virtus/internal/service"
)

type PoolHandler struct {
	pool *service.PoolService
}

func NewPoolHandler(pool *service.PoolService) *PoolHandler {
	return &PoolHandler{pool: pool}
}

func (h *PoolHandler) Status(w http.ResponseWriter, r *http.Request) {
	p, err := h.pool.GetStatus(r.Context())
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, p)
}
