package service

import (
	"context"
	"fmt"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/eulerbutcooler/virtus/pkg/crypto"
	"github.com/eulerbutcooler/virtus/pkg/pagination"
	"github.com/google/uuid"
)

type UpdateProfileInput struct {
	Name *string
}

type ChangePasswordInput struct {
	CurrentPassword string
	NewPassword     string
}

type UserService struct {
	users domain.UserRepository
}

func NewUserService(users domain.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("userService.GetByID: %w", err)
	}
	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, id uuid.UUID, in UpdateProfileInput) (*domain.User, error) {
	if in.Name != nil && len(*in.Name) == 0 {
		return nil, fmt.Errorf("%w: name cannot be empty", domain.ErrInvalidInput)
	}
	user, err := s.users.Update(ctx, id, domain.UpdateUserParams{
		Name: in.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("userService.UpdateProfile: %w", err)
	}
	return user, nil
}

func (s *UserService) ChangePassword(ctx context.Context, id uuid.UUID, in ChangePasswordInput) error {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("userService.ChangePassword: fetch user: %w", err)
	}
	if err := crypto.CheckPassword(user.PasswordHash, in.CurrentPassword); err != nil {
		return domain.ErrUnauthorized
	}
	if in.CurrentPassword == in.NewPassword {
		return fmt.Errorf("%w: new password must differ from current password", domain.ErrInvalidInput)
	}
	hash, err := crypto.HashPassword(in.NewPassword)
	if err != nil {
		return fmt.Errorf("userService.ChangePassword: hash new password: %w", err)
	}
	if err := s.users.UpdatePasswordHash(ctx, id, hash); err != nil {
		return fmt.Errorf("userService.ChangePassword: persist: %w", err)
	}
	return nil
}

func (s *UserService) VerifyUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	verified := true
	user, err := s.users.Update(ctx, id, domain.UpdateUserParams{
		Verified: &verified,
	})
	if err != nil {
		return nil, fmt.Errorf("userService.VerifyUser: %w", err)
	}
	return user, nil
}

func (s *UserService) List(ctx context.Context, page pagination.Page) ([]domain.User, pagination.Page, error) {
	users, total, err := s.users.List(ctx, page.Limit, page.Offset())
	if err != nil {
		return nil, pagination.Page{}, fmt.Errorf("userService.List: %w", err)
	}
	return users, page.WithTotal(total), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.users.Delete(ctx, id); err != nil {
		return fmt.Errorf("userService.DeleteUser: %w", err)
	}
	return nil
}
