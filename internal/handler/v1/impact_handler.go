package v1

import (
	"encoding/json"
	"net/http"

	"github.com/eulerbutcooler/virtus/internal/handler/middleware"
	"github.com/eulerbutcooler/virtus/internal/handler/response"
	"github.com/eulerbutcooler/virtus/internal/service"
	"github.com/eulerbutcooler/virtus/pkg/pagination"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ImpactHandler struct {
	impact *service.ImpactService
}

func NewImpactHandler(impact *service.ImpactService) *ImpactHandler {
	return &ImpactHandler{impact: impact}
}

// Member: submit a follow-up impact record after delivery.
// POST /impact
type recordImpactBody struct {
	DeliveryID         uuid.UUID       `json:"delivery_id"          validate:"required"`
	IntervalLabel      string          `json:"interval_label"       validate:"required,max=50"`
	OutcomeDescription *string         `json:"outcome_description"`
	SatisfactionScore  *int            `json:"satisfaction_score"   validate:"omitempty,min=1,max=5"`
	Metrics            json.RawMessage `json:"metrics"`
}

func (h *ImpactHandler) Record(w http.ResponseWriter, r *http.Request) {
	var body recordImpactBody
	if err := response.DecodeAndValidate(r, &body); err != nil {
		response.FromError(w, err)
		return
	}
	record, err := h.impact.Record(r.Context(), service.RecordInput{
		DeliveryID:         body.DeliveryID,
		UserID:             middleware.UserIDFrom(r.Context()),
		IntervalLabel:      body.IntervalLabel,
		OutcomeDescription: body.OutcomeDescription,
		SatisfactionScore:  body.SatisfactionScore,
		Metrics:            body.Metrics,
	})
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusCreated, record)
}

// GET /impact/{id}
func (h *ImpactHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid impact record id")
		return
	}
	record, err := h.impact.GetByID(r.Context(), id)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, record)
}

// GET /impact/delivery/{deliveryID}
func (h *ImpactHandler) ListByDelivery(w http.ResponseWriter, r *http.Request) {
	deliveryID, err := uuid.Parse(chi.URLParam(r, "deliveryID"))
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid delivery id")
		return
	}
	records, err := h.impact.ListByDelivery(r.Context(), deliveryID)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, records)
}

// Member: list my own impact records (paginated).
// GET /impact
func (h *ImpactHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	page := pagination.FromRequest(r)
	records, page, err := h.impact.ListByUser(r.Context(), middleware.UserIDFrom(r.Context()), page)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.List(w, http.StatusOK, records, page)
}
