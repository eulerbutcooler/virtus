package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/eulerbutcooler/virtus/pkg/pagination"
	"github.com/google/uuid"
)

type InstitutionService struct {
	institutions domain.InstitutionRepository
	pool         *PoolService
}

func NewInstitutionService(institutions domain.InstitutionRepository, pool *PoolService) *InstitutionService {
	return &InstitutionService{institutions: institutions, pool: pool}
}

type CreateInstitutionInput struct {
	UserID       uuid.UUID
	Name         string
	Type         domain.InstitutionType
	ContactEmail string
	Website      *string
	ESGGoals     json.RawMessage
}

type UpdateInstitutionInput struct {
	Name         *string
	Type         *domain.InstitutionType
	ContactEmail *string
	Website      *string
	ESGGoals     json.RawMessage
}

// Registers a new institution for the given user.
// Returns ErrConflict if the user already has one.
func (s *InstitutionService) Create(ctx context.Context, in CreateInstitutionInput) (*domain.Institution, error) {
	existing, err := s.institutions.GetByUserID(ctx, in.UserID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("institutionService.Create: check existing: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: user already has a registered institution", domain.ErrConflict)
	}
	inst, err := s.institutions.Create(ctx, domain.CreateInstitutionParams{
		UserID:       in.UserID,
		Name:         in.Name,
		Type:         in.Type,
		ContactEmail: in.ContactEmail,
		Website:      in.Website,
		ESGGoals:     in.ESGGoals,
	})
	if err != nil {
		return nil, fmt.Errorf("institutionService.Create: %w", err)
	}
	return inst, nil
}

func (s *InstitutionService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Institution, error) {
	inst, err := s.institutions.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("institutionService.GetByID: %w", err)
	}
	return inst, nil
}

// Returns the institution owned by the given user.
func (s *InstitutionService) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Institution, error) {
	inst, err := s.institutions.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("institutionService.GetByUserID: %w", err)
	}
	return inst, nil
}

// Applies partial edits. Only the institution owner or an admin may call this.
func (s *InstitutionService) Update(ctx context.Context, id uuid.UUID, requesterID uuid.UUID, requesterRole domain.UserRole, in UpdateInstitutionInput) (*domain.Institution, error) {
	inst, err := s.institutions.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("institutionService.Update: fetch: %w", err)
	}
	if inst.UserID != requesterID && requesterRole != domain.RoleAdmin {
		return nil, domain.ErrForbidden
	}
	updated, err := s.institutions.Update(ctx, id, domain.UpdateInstitutionParams{
		Name:         in.Name,
		Type:         in.Type,
		ContactEmail: in.ContactEmail,
		Website:      in.Website,
		ESGGoals:     in.ESGGoals,
	})
	if err != nil {
		return nil, fmt.Errorf("institutionService.Update: %w", err)
	}
	return updated, nil
}

// Marks the institution as verified. Admin only.
func (s *InstitutionService) Verify(ctx context.Context, id uuid.UUID) (*domain.Institution, error) {
	verified := true
	inst, err := s.institutions.Update(ctx, id, domain.UpdateInstitutionParams{Verified: &verified})
	if err != nil {
		return nil, fmt.Errorf("institutionService.Verify: %w", err)
	}
	return inst, nil
}

// Records an institutional contribution and credits the pool immediately.
// Institutional contributions are treated as already-completed (no payment intent flow).
func (s *InstitutionService) Contribute(ctx context.Context, institutionID uuid.UUID, amount float64, currency string, categoryTag, regionTag *string) (*domain.InstitutionalContribution, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("%w: amount must be positive", domain.ErrInvalidInput)
	}
	contrib, err := s.institutions.CreateContribution(ctx, institutionID, amount, currency, categoryTag, regionTag)
	if err != nil {
		return nil, fmt.Errorf("institutionService.Contribute: %w", err)
	}
	if err := s.pool.Credit(ctx, domain.GlobalPoolID, amount); err != nil {
		return nil, fmt.Errorf("institutionService.Contribute: credit pool: %w", err)
	}
	return contrib, nil
}

// Returns a paginated list of contributions for an institution.
func (s *InstitutionService) ListContributions(ctx context.Context, institutionID uuid.UUID, page pagination.Page) ([]*domain.InstitutionalContribution, pagination.Page, error) {
	items, total, err := s.institutions.ListContributions(ctx, institutionID, page.Limit, page.Offset())
	if err != nil {
		return nil, pagination.Page{}, fmt.Errorf("institutionService.ListContributions: %w", err)
	}
	return items, page.WithTotal(total), nil
}
