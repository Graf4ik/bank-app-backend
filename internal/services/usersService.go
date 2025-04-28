package services

import (
	"bank-app-backend/internal/entities"
	"bank-app-backend/internal/repository"
	"context"
	"fmt"
)

type UsersService interface {
	Me(ctx context.Context, userID uint) (*entities.User, error)
	GetAll(ctx context.Context) ([]*entities.User, error)
	Update(ctx context.Context, userID uint, input *entities.UpdateUserRequest) (*entities.User, error)
}

type usersService struct {
	repo repository.UsersRepository
}

func NewUsersService(r repository.UsersRepository) UsersService {
	return &usersService{repo: r}
}

func (s *usersService) Me(ctx context.Context, userID uint) (*entities.User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	return user, nil
}

func (s *usersService) GetAll(ctx context.Context) ([]*entities.User, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("users not found: %v", err)
	}

	return users, nil
}

func (s *usersService) Update(ctx context.Context, userID uint, input *entities.UpdateUserRequest) (*entities.User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Username != nil {
		user.Username = *input.Username
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return user, nil
}
