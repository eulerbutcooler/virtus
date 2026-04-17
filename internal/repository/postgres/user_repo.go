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

type userRepo struct {
	q *dbgen.Queries
}

func NewUserRepository(q *dbgen.Queries) domain.UserRepository {
	return &userRepo{q: q}
}

func (r *userRepo) Create(ctx context.Context, p domain.CreateUserParams) (*domain.User, error) {
	row, err := r.q.CreateUser(ctx, dbgen.CreateUserParams{
		Email:        p.Email,
		Name:         p.Name,
		PasswordHash: p.Password,
		Role:         dbgen.UserRoles(p.Role),
	})
	if err != nil {
		return nil, fmt.Errorf("userRepo.Create: %w", err)
	}
	return rowToUser(row), nil
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("userRepo.GetByID: %w", err)
	}
	return rowToUser(row), nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("userRepo.GetByEmail: %w", err)
	}
	return rowToUser(row), nil
}

func (r *userRepo) Update(ctx context.Context, id uuid.UUID, p domain.UpdateUserParams) (*domain.User, error) {
	row, err := r.q.UpdateUser(ctx, dbgen.UpdateUserParams{
		ID:       id,
		Name:     p.Name,
		Verified: p.Verified,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("userRepo.Update: %w", err)
	}
	return rowToUser(row), nil
}

func (r *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	n, err := r.q.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("userRepo.Delete: %w", err)
	}
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *userRepo) List(ctx context.Context, limit, offset int) ([]domain.User, int, error) {
	rows, err := r.q.ListUsers(ctx, dbgen.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("userRepo.List: %w", err)
	}
	total, err := r.q.CountUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("userRepo.List count: %w", err)
	}
	users := make([]domain.User, len(rows))
	for i, row := range rows {
		users[i] = *rowToUser(row)
	}
	return users, int(total), nil
}

// Maps the sqlc-generated User to the domain User.
func rowToUser(row dbgen.User) *domain.User {
	return &domain.User{
		ID:           row.ID,
		Email:        row.Email,
		Name:         row.Name,
		PasswordHash: row.PasswordHash,
		Role:         domain.UserRole(row.Role),
		Verified:     row.Verified,
		JoinedAt:     row.JoinedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}
