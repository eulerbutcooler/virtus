package v1

import (
	"net/http"

	"github.com/eulerbutcooler/virtus/internal/handler/response"
	"github.com/eulerbutcooler/virtus/internal/service"
	"github.com/eulerbutcooler/virtus/pkg/pagination"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type QueueHandler struct {
	queue *service.QueueService
}

func NewQueueHandler(queue *service.QueueService) *QueueHandler {
	return &QueueHandler{queue: queue}
}

func (h *QueueHandler) List(w http.ResponseWriter, r *http.Request) {
	page := pagination.FromRequest(r)
	entries, page, err := h.queue.ListAll(r.Context(), page)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.List(w, http.StatusOK, entries, page)
}

func (h *QueueHandler) GetByRequest(w http.ResponseWriter, r *http.Request) {
	requestID, err := uuid.Parse(chi.URLParam(r, "requestID"))
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid request id")
		return
	}
	entry, err := h.queue.GetByRequestID(r.Context(), requestID)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, entry)
}
