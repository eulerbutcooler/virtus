package v1

import (
	"net/http"
	"time"

	"github.com/eulerbutcooler/virtus/internal/handler/middleware"
	"github.com/eulerbutcooler/virtus/internal/handler/response"
	"github.com/eulerbutcooler/virtus/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type DeliveryHandler struct {
	deliveries *service.DeliveryService
}

func NewDeliveryHandler(deliveries *service.DeliveryService) *DeliveryHandler {
	return &DeliveryHandler{deliveries: deliveries}
}

// Admin: create a delivery record when an item ships.
// POST /admin/deliveries
type shipBody struct {
	FulfillmentID  uuid.UUID `json:"fulfillment_id"  validate:"required"`
	TrackingNumber *string   `json:"tracking_number"`
	Carrier        *string   `json:"carrier"`
}

func (h *DeliveryHandler) Ship(w http.ResponseWriter, r *http.Request) {
	var body shipBody
	if err := response.DecodeAndValidate(r, &body); err != nil {
		response.FromError(w, err)
		return
	}
	d, err := h.deliveries.Ship(r.Context(), service.ShipInput{
		FulfillmentID:  body.FulfillmentID,
		TrackingNumber: body.TrackingNumber,
		Carrier:        body.Carrier,
	})
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusCreated, d)
}

// GET /deliveries/{id}
func (h *DeliveryHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid delivery id")
		return
	}
	d, err := h.deliveries.GetByID(r.Context(), id)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, d)
}

// GET /deliveries/fulfillment/{fulfillmentID}
func (h *DeliveryHandler) GetByFulfillment(w http.ResponseWriter, r *http.Request) {
	fulfillmentID, err := uuid.Parse(chi.URLParam(r, "fulfillmentID"))
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid fulfillment id")
		return
	}
	d, err := h.deliveries.GetByFulfillmentID(r.Context(), fulfillmentID)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, d)
}

// Admin: verify a delivery with photo proof.
// POST /admin/deliveries/{id}/verify
type verifyDeliveryBody struct {
	ProofPhotoURL string    `json:"proof_photo_url" validate:"required,url"`
	DeliveredAt   time.Time `json:"delivered_at"    validate:"required"`
}

func (h *DeliveryHandler) Verify(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid delivery id")
		return
	}
	var body verifyDeliveryBody
	if err := response.DecodeAndValidate(r, &body); err != nil {
		response.FromError(w, err)
		return
	}
	d, err := h.deliveries.Verify(r.Context(), id, service.VerifyInput{
		ProofPhotoURL: body.ProofPhotoURL,
		DeliveredAt:   body.DeliveredAt,
		VerifiedBy:    middleware.UserIDFrom(r.Context()),
	})
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, d)
}

// Admin: mark a delivery as failed (lost/returned).
// POST /admin/deliveries/{id}/fail
func (h *DeliveryHandler) MarkFailed(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid delivery id")
		return
	}
	if err := h.deliveries.MarkFailed(r.Context(), id); err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, map[string]string{"status": "failed"})
}
