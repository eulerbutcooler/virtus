package v1

import (
	"io"
	"net/http"

	"github.com/eulerbutcooler/virtus/internal/handler/middleware"
	"github.com/eulerbutcooler/virtus/internal/handler/response"
	"github.com/eulerbutcooler/virtus/internal/service"
	"github.com/eulerbutcooler/virtus/pkg/pagination"
	stripepkg "github.com/eulerbutcooler/virtus/pkg/stripe"
)

// Can be used to swap out the Stripe provider in tests.
type webhookParser interface {
	ParseWebhook(payload []byte, sigHeader string) (*stripepkg.WebhookEvent, error)
}

type ContributionHandler struct {
	contributions *service.ContributionService
	webhook       webhookParser
}

func NewContributionHandler(contributions *service.ContributionService, webhook webhookParser) *ContributionHandler {
	return &ContributionHandler{contributions: contributions, webhook: webhook}
}

// Member creates a Stripe PaymentIntent and a pending contribution row.
// The client secret is returned so the frontend can confirm the payment.
// POST /contributions
type initiateBody struct {
	Amount   float64 `json:"amount"   validate:"required,gt=0"`
	Currency string  `json:"currency" validate:"required,len=3"`
}

func (h *ContributionHandler) Initiate(w http.ResponseWriter, r *http.Request) {
	var body initiateBody
	if err := response.DecodeAndValidate(r, &body); err != nil {
		response.FromError(w, err)
		return
	}
	result, err := h.contributions.Initiate(r.Context(), service.InitialInput{
		UserID:   middleware.UserIDFrom(r.Context()),
		Amount:   body.Amount,
		Currency: body.Currency,
	})
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusCreated, map[string]any{
		"contribution":  result.Contribution,
		"client_secret": result.ClientSecret,
	})
}

// Lists member's contributions (paginated).
// GET /contributions
func (h *ContributionHandler) List(w http.ResponseWriter, r *http.Request) {
	page := pagination.FromRequest(r)
	items, page, err := h.contributions.ListByUser(r.Context(), middleware.UserIDFrom(r.Context()), page)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.List(w, http.StatusOK, items, page)
}

// Gets a single contribution of member.
// GET /contributions/{id}
func (h *ContributionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid contribution id")
		return
	}
	c, err := h.contributions.GetByID(r.Context(), id)
	if err != nil {
		response.FromError(w, err)
		return
	}
	if !canAccess(r.Context(), c.UserID) {
		response.Fail(w, http.StatusForbidden, "forbidden")
		return
	}
	response.OK(w, http.StatusOK, c)
}

// Total amount contributed by member
// GET /contributions/total
func (h *ContributionHandler) Total(w http.ResponseWriter, r *http.Request) {
	total, err := h.contributions.GetUserTotal(r.Context(), middleware.UserIDFrom(r.Context()))
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, map[string]float64{"total": total})
}

// Public: receives Stripe webhook events.
// POST /webhooks/stripe  (no auth middleware — must be outside the RequireAuth group)
// Stripe may deliver the same event more than once; the underlying service
// methods are idempotent so duplicate deliveries are safe.
func (h *ContributionHandler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1 MB cap
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "failed to read request body")
		return
	}
	sig := r.Header.Get("Stripe-Signature")
	if sig == "" {
		response.Fail(w, http.StatusBadRequest, "missing Stripe-Signature header")
		return
	}
	ev, err := h.webhook.ParseWebhook(payload, sig)
	if err != nil {
		// Return 400 so Stripe knows the event was not processed and will retry.
		response.Fail(w, http.StatusBadRequest, "invalid webhook signature")
		return
	}

	switch ev.Type {
	case "payment_intent.succeeded":
		if err := h.contributions.Complete(r.Context(), ev.PaymentIntentID, ev.ChargeID); err != nil {
			response.FromError(w, err)
			return
		}
	case "payment_intent.payment_failed", "payment_intent.canceled":
		if err := h.contributions.Fail(r.Context(), ev.PaymentIntentID); err != nil {
			response.FromError(w, err)
			return
		}
	}

	// Return 200 for all handled and unhandled event types
	// Stripe won't retry on 2xx
	response.OK(w, http.StatusOK, map[string]string{"received": "ok"})
}
