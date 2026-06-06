package v1

import (
	"net/http"
	"time"

	"github.com/eulerbutcooler/virtus/internal/handler/response"
	"github.com/eulerbutcooler/virtus/internal/service"
	"github.com/eulerbutcooler/virtus/pkg/pagination"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FulfillmentHandler struct {
	fulfillments *service.FulfillmentService
}

func NewFulfillmentHandler(fulfillments *service.FulfillmentService) *FulfillmentHandler {
	return &FulfillmentHandler{fulfillments: fulfillments}
}

// Admin: start procurement for a funded request.
// POST /admin/fulfillments
type beginFulfillmentBody struct {
	RequestID uuid.UUID `json:"request_id" validate:"required"`
}

func (h *FulfillmentHandler) Begin(w http.ResponseWriter, r *http.Request) {
	var body beginFulfillmentBody
	if err := response.DecodeAndValidate(r, &body); err != nil {
		response.FromError(w, err)
		return
	}
	f, err := h.fulfillments.Begin(r.Context(), body.RequestID)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusCreated, f)
}

// Admin: list all fulfillments (paginated).
// GET /admin/fulfillments
func (h *FulfillmentHandler) List(w http.ResponseWriter, r *http.Request) {
	page := pagination.FromRequest(r)
	items, page, err := h.fulfillments.List(r.Context(), page)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.List(w, http.StatusOK, items, page)
}

// Admin/member: get a fulfillment by its ID.
// GET /fulfillments/{id}
func (h *FulfillmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid fulfillment id")
		return
	}
	f, err := h.fulfillments.GetByID(r.Context(), id)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, f)
}

// Member: get the fulfillment for a specific request.
// GET /fulfillments/request/{requestID}
func (h *FulfillmentHandler) GetByRequest(w http.ResponseWriter, r *http.Request) {
	requestID, err := uuid.Parse(chi.URLParam(r, "requestID"))
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid request id")
		return
	}
	f, err := h.fulfillments.GetByRequestID(r.Context(), requestID)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, f)
}

// Admin: update vendor info, cost, status, notes.
// PATCH /admin/fulfillments/{id}
type updateFulfillmentBody struct {
	VendorName *string    `json:"vendor_name" validate:"omitempty,min=1"`
	VendorRef  *string    `json:"vendor_ref"  validate:"omitempty,min=1"`
	ActualCost *float64   `json:"actual_cost" validate:"omitempty,gt=0"`
	Status     *string    `json:"status"      validate:"omitempty,oneof=pending vendor_selected ordered shipped delivered cancelled"`
	Notes      *string    `json:"notes"       validate:"omitempty,min=1"`
	ProcuredAt *time.Time `json:"procured_at"`
}

func (h *FulfillmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid fulfillment id")
		return
	}
	var body updateFulfillmentBody
	if err := response.DecodeAndValidate(r, &body); err != nil {
		response.FromError(w, err)
		return
	}
	in := service.UpdateDetailsInput{
		VendorName: body.VendorName,
		VendorRef:  body.VendorRef,
		ActualCost: body.ActualCost,
		Notes:      body.Notes,
		ProcuredAt: body.ProcuredAt,
	}
	if body.Status != nil {
		s := parseFulfillmentStatus(*body.Status)
		in.Status = &s
	}
	f, err := h.fulfillments.UpdateDetails(r.Context(), id, in)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, f)
}

// Admin: cancel a fulfillment.
// POST /admin/fulfillments/{id}/cancel
func (h *FulfillmentHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid fulfillment id")
		return
	}
	f, err := h.fulfillments.Cancel(r.Context(), id)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, f)
}
