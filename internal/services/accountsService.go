package services

import (
	"bank-app-backend/internal/entities"
	"bank-app-backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type AccountsService interface {
	GetAll(ctx context.Context, userID uint) ([]*entities.Account, error)
	GetByID(ctx context.Context, userID, accountID uint) (*entities.Account, error)
	Create(ctx context.Context, userID uint, input *entities.CreateAccountRequest) (*entities.Account, error)
	Delete(ctx context.Context, userID, accountID uint) error
}

type accountsService struct {
	repo repository.AccountsRepository
}

func NewAccountsService(r repository.AccountsRepository) AccountsService {
	return &accountsService{repo: r}
}

func (s accountsService) GetAll(ctx context.Context, userID uint) ([]*entities.Account, error) {
	accounts, err := s.repo.GetAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %v", err)
	}

	return accounts, nil
}

func (s accountsService) GetByID(ctx context.Context, userID, accountID uint) (*entities.Account, error) {
	account, err := s.repo.GetByID(ctx, userID, accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return account, nil
}

func (s accountsService) Create(ctx context.Context, userID uint, req *entities.CreateAccountRequest) (*entities.Account, error) {
	account := &entities.Account{
		UserID:   userID,
		Type:     req.Type,
		Currency: req.Currency,
		Balance:  0,
		Status:   "active",
	}

	if err := s.repo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

func (s accountsService) Delete(ctx context.Context, userID, accountID uint) error {
	account, err := s.repo.GetByID(ctx, userID, accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("account not found")
		}
		return fmt.Errorf("failed to fetch account: %w", err)
	}

	if account.Balance != 0 {
		return fmt.Errorf("cannot close account with non-zero balance")
	}

	account.Status = "closed"

	if err := s.repo.Update(ctx, account); err != nil {
		return fmt.Errorf("failed to close account: %w", err)
	}

	return nil
}
