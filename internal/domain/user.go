package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleMember      UserRole = "member"
	RoleAdmin       UserRole = "admin"
	RoleInstitution UserRole = "institution"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	Verified     bool      `json:"verified"`
	JoinedAt     time.Time `json:"joined_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateUserParams struct {
	Email    string
	Name     string
	Password string
	Role     UserRole
}

type UpdateUserParams struct {
	Name     *string
	Verified *bool
}

type UserRepository interface {
	Create(ctx context.Context, p CreateUserParams) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, id uuid.UUID, p UpdateUserParams) (*User, error)
	UpdatePasswordHash(ctx context.Context, id uuid.UUID, hash string) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]User, int, error)
}
