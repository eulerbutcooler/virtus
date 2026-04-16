package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type FulfillmentStatus string

const (
	FulfillmentPending        FulfillmentStatus = "pending"
	FulfillmentVendorSelected FulfillmentStatus = "vendor_selected"
	FulfillmentOrdered        FulfillmentStatus = "ordered"
	FulfillmentShipped        FulfillmentStatus = "shipped"
	FulfillmentDelivered      FulfillmentStatus = "delivered"
	FulfillmentCancelled      FulfillmentStatus = "cancelled"
)

type Fulfillment struct {
	ID                uuid.UUID         `json:"id"`
	RequestID         uuid.UUID         `json:"request_id"`
	VendorName        *string           `json:"vendor_name,omitempty"`
	VendorRef         *string           `json:"vendor_ref,omitempty"`
	ActualCost        *float64          `json:"actual_cost,omitempty"`
	ProcurementStatus FulfillmentStatus `json:"procurement_status"`
	Notes             *string           `json:"notes,omitempty"`
	ProcuredAt        *time.Time        `json:"procured_at,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

type CreateFulfillmentParams struct {
	RequestID uuid.UUID
}

type UpdateFulfillmentParams struct {
	VendorName *string
	VendorRef  *string
	ActualCost *float64
	Status     *FulfillmentStatus
	Notes      *string
	ProcuredAt *time.Time
}

type FulfillmentRepository interface {
	Create(ctx context.Context, p CreateFulfillmentParams) (*Fulfillment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Fulfillment, error)
	GetByRequestID(ctx context.Context, requestID uuid.UUID) (*Fulfillment, error)
	Update(ctx context.Context, id uuid.UUID, p UpdateFulfillmentParams) (*Fulfillment, error)
	List(ctx context.Context, limit, offset int) ([]*Fulfillment, int, error)
}
