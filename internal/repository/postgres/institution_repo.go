package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/eulerbutcooler/virtus/internal/domain"
	dbgen "github.com/eulerbutcooler/virtus/internal/repository/postgres/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type institutionRepo struct {
	q *dbgen.Queries
}

func NewInstitutionRepository(q *dbgen.Queries) domain.InstitutionRepository {
	return &institutionRepo{q: q}
}

func (r *institutionRepo) Create(ctx context.Context, p domain.CreateInstitutionParams) (*domain.Institution, error) {
	row, err := r.q.CreateInstitution(ctx, dbgen.CreateInstitutionParams{
		UserID:       p.UserID,
		Name:         p.Name,
		Type:         string(p.Type),
		ContactEmail: p.ContactEmail,
		Website:      p.Website,
		EsgGoals:     p.ESGGoals,
	})
	if err != nil {
		return nil, fmt.Errorf("institutionRepo.Create: %w", err)
	}
	return rowToInstitution(row), nil
}

func (r *institutionRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Institution, error) {
	row, err := r.q.GetInstitutionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("institutionRepo.GetByID: %w", err)
	}
	return rowToInstitution(row), nil
}

func (r *institutionRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Institution, error) {
	row, err := r.q.GetInstitutionByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("institutionRepo.GetByUserID: %w", err)
	}
	return rowToInstitution(row), nil
}

func (r *institutionRepo) Update(ctx context.Context, id uuid.UUID, p domain.UpdateInstitutionParams) (*domain.Institution, error) {
	var instType dbgen.NullInstitutionType
	if p.Type != nil {
		instType = dbgen.NullInstitutionType{
			InstitutionType: dbgen.InstitutionType(*p.Type),
			Valid:           true,
		}
	}

	row, err := r.q.UpdateInstitution(ctx, dbgen.UpdateInstitutionParams{
		ID:           id,
		Name:         p.Name,
		Type:         instType,
		ContactEmail: p.ContactEmail,
		Website:      p.Website,
		EsgGoals:     []byte(p.ESGGoals),
		Verified:     p.Verified,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("institutionRepo.Update: %w", err)
	}
	return rowToInstitution(row), nil
}

func (r *institutionRepo) CreateContribution(
	ctx context.Context,
	institutionID uuid.UUID,
	amount float64,
	currency string,
	categoryTag, regionTag *string,
) (*domain.InstitutionalContribution, error) {
	row, err := r.q.CreateInstitutionalContribution(ctx, dbgen.CreateInstitutionalContributionParams{
		InstitutionID: institutionID,
		PoolID:        domain.GlobalPoolID,
		Amount:        amount,
		Currency:      currency,
		CategoryTag:   categoryTag,
		RegionTag:     regionTag,
	})
	if err != nil {
		return nil, fmt.Errorf("institutionRepo.CreateContribution: %w", err)
	}
	return rowToInstitutionalContribution(row), nil
}

func (r *institutionRepo) ListContributions(ctx context.Context, institutionID uuid.UUID, limit, offset int) ([]*domain.InstitutionalContribution, int, error) {
	rows, err := r.q.ListInstitutionalContributions(ctx, dbgen.ListInstitutionalContributionsParams{
		InstitutionID: institutionID,
		Limit:         int32(limit),
		Offset:        int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("institutionRepo.ListContributions: %w", err)
	}
	total, err := r.q.CountInstitutionalContributions(ctx, institutionID)
	if err != nil {
		return nil, 0, fmt.Errorf("institutionRepo.ListContributions count: %w", err)
	}
	out := make([]*domain.InstitutionalContribution, len(rows))
	for i, row := range rows {
		out[i] = rowToInstitutionalContribution(row)
	}
	return out, int(total), nil
}

func rowToInstitution(row dbgen.Institution) *domain.Institution {
	return &domain.Institution{
		ID:           row.ID,
		UserID:       row.UserID,
		Name:         row.Name,
		Type:         domain.InstitutionType(row.Type),
		ContactEmail: row.ContactEmail,
		Website:      row.Website,
		ESGGoals:     row.EsgGoals,
		Verified:     row.Verified,
		JoinedAt:     row.JoinedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}

func rowToInstitutionalContribution(row dbgen.InstitutionalContribution) *domain.InstitutionalContribution {
	return &domain.InstitutionalContribution{
		ID:            row.ID,
		InstitutionID: row.InstitutionID,
		PoolID:        row.PoolID,
		Amount:        row.Amount,
		Currency:      row.Currency,
		Status:        domain.ContributionStatus(row.Status),
		PaymentRef:    row.PaymentRef,
		CategoryTag:   row.CategoryTag,
		RegionTag:     row.RegionTag,
		Notes:         row.Notes,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}
}
