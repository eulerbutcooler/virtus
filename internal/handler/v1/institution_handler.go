package v1

import (
	"encoding/json"
	"net/http"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/eulerbutcooler/virtus/internal/handler/middleware"
	"github.com/eulerbutcooler/virtus/internal/handler/response"
	"github.com/eulerbutcooler/virtus/internal/service"
	"github.com/eulerbutcooler/virtus/pkg/pagination"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type InstitutionHandler struct {
	institutions *service.InstitutionService
}

func NewInstitutionHandler(institutions *service.InstitutionService) *InstitutionHandler {
	return &InstitutionHandler{institutions: institutions}
}

// POST /institutions
type createInstitutionBody struct {
	Name         string                 `json:"name"          validate:"required,min=2,max=200"`
	Type         domain.InstitutionType `json:"type"          validate:"required,oneof=corporation ngo government foundation other"`
	ContactEmail string                 `json:"contact_email" validate:"required,email"`
	Website      *string                `json:"website"       validate:"omitempty,url"`
	ESGGoals     json.RawMessage        `json:"esg_goals"`
}

func (h *InstitutionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body createInstitutionBody
	if err := response.DecodeAndValidate(r, &body); err != nil {
		response.FromError(w, err)
		return
	}
	inst, err := h.institutions.Create(r.Context(), service.CreateInstitutionInput{
		UserID:       middleware.UserIDFrom(r.Context()),
		Name:         body.Name,
		Type:         body.Type,
		ContactEmail: body.ContactEmail,
		Website:      body.Website,
		ESGGoals:     body.ESGGoals,
	})
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusCreated, inst)
}

// GET /institutions/me
func (h *InstitutionHandler) GetMine(w http.ResponseWriter, r *http.Request) {
	inst, err := h.institutions.GetByUserID(r.Context(), middleware.UserIDFrom(r.Context()))
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, inst)
}

// GET /institutions/{id}
func (h *InstitutionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid institution id")
		return
	}
	inst, err := h.institutions.GetByID(r.Context(), id)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, inst)
}

// PATCH /institutions/{id}
type updateInstitutionBody struct {
	Name         *string                 `json:"name"          validate:"omitempty,min=2,max=200"`
	Type         *domain.InstitutionType `json:"type"          validate:"omitempty,oneof=corporation ngo government foundation other"`
	ContactEmail *string                 `json:"contact_email" validate:"omitempty,email"`
	Website      *string                 `json:"website"       validate:"omitempty,url"`
	ESGGoals     json.RawMessage         `json:"esg_goals"`
}

func (h *InstitutionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid institution id")
		return
	}
	var body updateInstitutionBody
	if err := response.DecodeAndValidate(r, &body); err != nil {
		response.FromError(w, err)
		return
	}
	inst, err := h.institutions.Update(r.Context(), id,
		middleware.UserIDFrom(r.Context()),
		middleware.UserRoleFrom(r.Context()),
		service.UpdateInstitutionInput{
			Name:         body.Name,
			Type:         body.Type,
			ContactEmail: body.ContactEmail,
			Website:      body.Website,
			ESGGoals:     body.ESGGoals,
		},
	)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, inst)
}

// Admin: verify an institution.
// POST /admin/institutions/{id}/verify
func (h *InstitutionHandler) Verify(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid institution id")
		return
	}
	inst, err := h.institutions.Verify(r.Context(), id)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, inst)
}

// POST /institutions/{id}/contributions
type institutionContributeBody struct {
	Amount      float64 `json:"amount"       validate:"required,gt=0"`
	Currency    string  `json:"currency"     validate:"required,len=3"`
	CategoryTag *string `json:"category_tag" validate:"omitempty,max=100"`
	RegionTag   *string `json:"region_tag"   validate:"omitempty,max=100"`
}

func (h *InstitutionHandler) Contribute(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid institution id")
		return
	}
	var body institutionContributeBody
	if err := response.DecodeAndValidate(r, &body); err != nil {
		response.FromError(w, err)
		return
	}
	contrib, err := h.institutions.Contribute(r.Context(), id, body.Amount, body.Currency, body.CategoryTag, body.RegionTag)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusCreated, contrib)
}

// GET /institutions/{id}/contributions
func (h *InstitutionHandler) ListContributions(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid institution id")
		return
	}
	page := pagination.FromRequest(r)
	items, page, err := h.institutions.ListContributions(r.Context(), id, page)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.List(w, http.StatusOK, items, page)
}
