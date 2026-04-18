package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type InstitutionType string

const (
	InstitutionCorporation InstitutionType = "corporation"
	InstitutionNGO         InstitutionType = "ngo"
	InstitutionGovernment  InstitutionType = "government"
	InstitutionFoundation  InstitutionType = "foundation"
	InstitutionOther       InstitutionType = "other"
)

type Institution struct {
	ID           uuid.UUID       `json:"id"`
	UserID       uuid.UUID       `json:"user_id"`
	Name         string          `json:"name"`
	Type         InstitutionType `json:"type"`
	ContactEmail string          `json:"contact_email"`
	Website      *string         `json:"website,omitempty"`
	ESGGoals     json.RawMessage `json:"esg_goals"`
	Verified     bool            `json:"verified"`
	JoinedAt     time.Time       `json:"joined_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type InstitutionalContribution struct {
	ID            uuid.UUID          `json:"id"`
	InstitutionID uuid.UUID          `json:"institution_id"`
	PoolID        uuid.UUID          `json:"pool_id"`
	Amount        float64            `json:"amount"`
	Currency      string             `json:"currency"`
	Status        ContributionStatus `json:"status"`
	PaymentRef    *string            `json:"payment_ref,omitempty"`
	CategoryTag   *string            `json:"category_tag,omitempty"`
	RegionTag     *string            `json:"region_tag,omitempty"`
	Notes         *string            `json:"notes,omitempty"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

type CreateInstitutionParams struct {
	UserID       uuid.UUID
	Name         string
	Type         InstitutionType
	ContactEmail string
	Website      *string
	ESGGoals     json.RawMessage
}

// UpdateInstitutionParams uses pointer fields so only non-nil values are applied.
type UpdateInstitutionParams struct {
	Name         *string
	Type         *InstitutionType
	ContactEmail *string
	Website      *string
	ESGGoals     json.RawMessage // nil means no change
	Verified     *bool
}

type InstitutionRepository interface {
	Create(ctx context.Context, p CreateInstitutionParams) (*Institution, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Institution, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Institution, error)
	Update(ctx context.Context, id uuid.UUID, p UpdateInstitutionParams) (*Institution, error)
	CreateContribution(ctx context.Context, institutionID uuid.UUID, amount float64, currency string, categoryTag, regionTag *string) (*InstitutionalContribution, error)
	ListContributions(ctx context.Context, institutionID uuid.UUID, limit, offset int) ([]*InstitutionalContribution, int, error)
}
