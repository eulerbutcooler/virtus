// Package stripe wraps the Stripe Go SDK and exposes only what the application
// needs: a PaymentProvider for creating payment intents and a ParseWebhook
// helper for processing incoming webhook events.
package stripe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eulerbutcooler/virtus/internal/service"
	stripelib "github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/stripe/stripe-go/v81/webhook"
)

// Holds the fields the application needs from a Stripe event.
// Callers do not need to import stripe-go directly.
type WebhookEvent struct {
	Type            string
	PaymentIntentID string
	ChargeID        string
}

type Provider struct {
	webhookSecret string
}

// Sets the global Stripe secret key and returns a Provider.
func New(secretKey, webhookSecret string) *Provider {
	stripelib.Key = secretKey
	return &Provider{webhookSecret: webhookSecret}
}

// Implements service.PaymentProvider.
// Amount is in the currency's major unit (e.g. 10.00 for $10.00 USD).
// Stripe requires amounts in the smallest currency unit so we multiply by 100.
func (p *Provider) CreateIntent(ctx context.Context, amount float64, currency string, metadata map[string]string) (*service.PaymentIntent, error) {
	params := &stripelib.PaymentIntentParams{
		Amount:   stripelib.Int64(int64(amount * 100)),
		Currency: stripelib.String(currency),
		Metadata: metadata,
	}
	params.Context = ctx

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("stripe: create intent: %w", err)
	}
	return &service.PaymentIntent{
		IntentID:     pi.ID,
		ClientSecret: pi.ClientSecret,
	}, nil
}

// Verifies the Stripe-Signature header and extracts the fields we care about.
// Returns an error if the signature is invalid.
func (p *Provider) ParseWebhook(payload []byte, sigHeader string) (*WebhookEvent, error) {
	event, err := webhook.ConstructEvent(payload, sigHeader, p.webhookSecret)
	if err != nil {
		return nil, fmt.Errorf("stripe: verify signature: %w", err)
	}

	ev := &WebhookEvent{Type: string(event.Type)}

	switch event.Type {
	case "payment_intent.succeeded",
		"payment_intent.payment_failed",
		"payment_intent.canceled":
		var pi stripelib.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			return nil, fmt.Errorf("stripe: unmarshal PaymentIntent: %w", err)
		}
		ev.PaymentIntentID = pi.ID
		if pi.LatestCharge != nil {
			ev.ChargeID = pi.LatestCharge.ID
		}
	}

	return ev, nil
}
