package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RequestStatus string
type UrgencyLevel string

const (
	RequestDraft     RequestStatus = "draft"
	RequestSubmitted RequestStatus = "submitted"
	RequestVerified  RequestStatus = "verified"
	RequestQueued    RequestStatus = "queued"
	RequestFunded    RequestStatus = "funded"
	RequestProcuring RequestStatus = "procuring"
	RequestDelivered RequestStatus = "delivered"
	RequestCompleted RequestStatus = "completed"
	RequestRejected  RequestStatus = "rejected"
)

const (
	UrgencyCritical UrgencyLevel = "critical"
	UrgencyHigh     UrgencyLevel = "high"
	UrgencyStandard UrgencyLevel = "standard"
	UrgencyLow      UrgencyLevel = "low"
)

// Returns the scoring weight for a given urgency level.
func (u UrgencyLevel) Factor() float64 {
	switch u {
	case UrgencyCritical:
		return 1.0
	case UrgencyHigh:
		return 0.7
	case UrgencyStandard:
		return 0.4
	default:
		return 0.2
	}
}

type Request struct {
	ID            uuid.UUID     `json:"id"`
	UserID        uuid.UUID     `json:"user_id"`
	ItemCategory  string        `json:"item_category"`
	ItemName      string        `json:"item_name"`
	Description   string        `json:"description"`
	Urgency       UrgencyLevel  `json:"urgency"`
	EstimatedCost float64       `json:"estimated_cost"`
	Justification string        `json:"justification"`
	Status        RequestStatus `json:"status"`
	RejectionNote *string       `json:"rejection_note,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type CreateRequestParams struct {
	UserID        uuid.UUID
	ItemCategory  string
	ItemName      string
	Description   string
	Urgency       UrgencyLevel
	EstimatedCost float64
	Justification string
}

type UpdateRequestParams struct {
	ItemCategory  *string
	ItemName      *string
	Description   *string
	Urgency       *UrgencyLevel
	EstimatedCost *float64
	Justification *string
}

type RequestFilter struct {
	UserID *uuid.UUID
	Status *RequestStatus
	Limit  int
	Offset int
}

type RequestRepository interface {
	Create(ctx context.Context, p CreateRequestParams) (*Request, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Request, error)
	Update(ctx context.Context, id uuid.UUID, p UpdateRequestParams) (*Request, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status RequestStatus, note *string) error
	List(ctx context.Context, f RequestFilter) ([]*Request, int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
