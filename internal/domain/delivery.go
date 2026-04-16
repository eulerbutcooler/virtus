package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DeliveryStatus string

const (
	DeliveryInTransit DeliveryStatus = "in_transit"
	DeliveryDelivered DeliveryStatus = "delivered"
	DeliveryFailed    DeliveryStatus = "failed"
	DeliveryReturned  DeliveryStatus = "returned"
)

type Delivery struct {
	ID             uuid.UUID      `json:"id"`
	FulfillmentID  uuid.UUID      `json:"fulfillment_id"`
	TrackingNumber *string        `json:"tracking_number,omitempty"`
	Carrier        *string        `json:"carrier,omitempty"`
	ProofPhotoURL  *string        `json:"proof_photo_url,omitempty"`
	Status         DeliveryStatus `json:"status"`
	DeliveredAt    *time.Time     `json:"delivered_at,omitempty"`
	VerifiedAt     *time.Time     `json:"verified_at,omitempty"`
	VerifiedBy     *uuid.UUID     `json:"verified_by,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type CreateDeliveryParams struct {
	FulfillmentID  uuid.UUID
	TrackingNumber *string
	Carrier        *string
}

type VerifyDeliveryParams struct {
	ProofPhotoURL string
	DeliveredAt   time.Time
	VerifiedBy    uuid.UUID
}

type DeliveryRepository interface {
	Create(ctx context.Context, p CreateDeliveryParams) (*Delivery, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Delivery, error)
	GetByFulfillmentID(ctx context.Context, fulfillmentID uuid.UUID) (*Delivery, error)
	Verify(ctx context.Context, id uuid.UUID, p VerifyDeliveryParams) (*Delivery, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status DeliveryStatus) error
}
