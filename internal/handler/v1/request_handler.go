package v1

import (
	"net/http"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/eulerbutcooler/virtus/internal/handler/middleware"
	"github.com/eulerbutcooler/virtus/internal/handler/response"
	"github.com/eulerbutcooler/virtus/internal/service"
	"github.com/eulerbutcooler/virtus/pkg/pagination"
)

type RequestHandler struct {
	requests *service.RequestService
}

func NewRequestHandler(requests *service.RequestService) *RequestHandler {
	return &RequestHandler{requests: requests}
}

// Request bodies

type submitRequestBody struct {
	ItemCategory  string              `json:"item_category" validate:"required,max=100"`
	ItemName      string              `json:"item_name" validate:"required,max=200"`
	Description   string              `json:"description" validate:"required"`
	Urgency       domain.UrgencyLevel `json:"urgency" validate:"required,oneof=critical high standard low"`
	EstimatedCost float64             `json:"estimated_cost" validate:"required,gt=0"`
	Justification string              `json:"justification" validate:"required"`
}

type updateRequestBody struct {
	ItemCategory  *string              `json:"item_category" validate:"omitempty,max=100"`
	ItemName      *string              `json:"item_name" validate:"omitempty,max=200"`
	Description   *string              `json:"description" validate:"omitempty,min=1"`
	Urgency       *domain.UrgencyLevel `json:"urgency" validate:"omitempty,oneof=critical high standard low"`
	EstimatedCost *float64             `json:"estimated_cost" validate:"omitempty,gt=0"`
	Justification *string              `json:"justification" validate:"omitempty,min=1"`
}

type rejectRequestBody struct {
	Note string `json:"note" validate:"required,max=1000"`
}

// Member actions

func (h *RequestHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var body submitRequestBody
	if err := response.DecodeAndValidate(r, &body); err != nil {
		response.FromError(w, err)
		return
	}
	req, err := h.requests.Submit(r.Context(), service.SubmitInput{
		UserID:        middleware.UserIDFrom(r.Context()),
		ItemCategory:  body.ItemCategory,
		ItemName:      body.ItemName,
		Description:   body.Description,
		Urgency:       body.Urgency,
		EstimatedCost: body.EstimatedCost,
		Justification: body.Justification,
	})
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusCreated, req)
}

func (h *RequestHandler) List(w http.ResponseWriter, r *http.Request) {
	page := pagination.FromRequest(r)
	reqs, page, err := h.requests.MyRequests(r.Context(), middleware.UserIDFrom(r.Context()), page)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.List(w, http.StatusOK, reqs, page)
}

func (h *RequestHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid request id")
		return
	}
	req, err := h.requests.GetByID(r.Context(), id)
	if err != nil {
		response.FromError(w, err)
		return
	}
	if !canAccess(r.Context(), req.UserID) {
		response.Fail(w, http.StatusForbidden, "forbidden")
		return
	}
	response.OK(w, http.StatusOK, req)
}

func (h *RequestHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid request id")
		return
	}
	// Ownership check before mutating
	existing, err := h.requests.GetByID(r.Context(), id)
	if err != nil {
		response.FromError(w, err)
		return
	}
	if existing.UserID != middleware.UserIDFrom(r.Context()) {
		response.Fail(w, http.StatusForbidden, "forbidden")
		return
	}
	var body updateRequestBody
	if err := response.DecodeAndValidate(r, &body); err != nil {
		response.FromError(w, err)
		return
	}
	req, err := h.requests.Update(r.Context(), id, service.UpdateRequestInput{
		ItemCategory:  body.ItemCategory,
		ItemName:      body.ItemName,
		Description:   body.Description,
		Urgency:       body.Urgency,
		EstimatedCost: body.EstimatedCost,
		Justification: body.Justification,
	})
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, req)
}

func (h *RequestHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid request id")
		return
	}
	existing, err := h.requests.GetByID(r.Context(), id)
	if err != nil {
		response.FromError(w, err)
		return
	}
	if existing.UserID != middleware.UserIDFrom(r.Context()) {
		response.Fail(w, http.StatusForbidden, "forbidden")
		return
	}
	if err := h.requests.Delete(r.Context(), id); err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// Admin actions

func (h *RequestHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	page := pagination.FromRequest(r)

	var status *domain.RequestStatus
	if raw := r.URL.Query().Get("status"); raw != "" {
		s, ok := parseRequestStatus(raw)
		if !ok {
			response.Fail(w, http.StatusBadRequest, "invalid status filter")
			return
		}
		status = s
	}

	reqs, page, err := h.requests.AdminList(r.Context(), status, page)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.List(w, http.StatusOK, reqs, page)
}

func (h *RequestHandler) Verify(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid request id")
		return
	}
	req, err := h.requests.Verify(r.Context(), id)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, req)
}

func (h *RequestHandler) Reject(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid request id")
		return
	}
	var body rejectRequestBody
	if err := response.DecodeAndValidate(r, &body); err != nil {
		response.FromError(w, err)
		return
	}
	req, err := h.requests.Reject(r.Context(), id, body.Note)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, req)
}

// Validates a status query param against the known set.
func parseRequestStatus(s string) (*domain.RequestStatus, bool) {
	switch domain.RequestStatus(s) {
	case domain.RequestDraft, domain.RequestSubmitted, domain.RequestVerified,
		domain.RequestQueued, domain.RequestFunded, domain.RequestProcuring,
		domain.RequestDelivered, domain.RequestCompleted, domain.RequestRejected:
		rs := domain.RequestStatus(s)
		return &rs, true
	default:
		return nil, false
	}
}
