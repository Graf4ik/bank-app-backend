package services

import (
	"bank-app-backend/internal/entities"
	"bank-app-backend/internal/lib/kafka"
	lib "bank-app-backend/internal/lib/logger"
	"bank-app-backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type AccountsService interface {
	GetAll(ctx context.Context, userID uint) ([]*entities.Account, error)
	GetByID(ctx context.Context, userID, accountID uint) (*entities.Account, error)
	Deposit(ctx context.Context, userID, accountID uint, amount float64) (*entities.Account, error)
	Create(ctx context.Context, userID uint, input *entities.CreateAccountRequest) (*entities.Account, error)
	Delete(ctx context.Context, userID, accountID uint) error
}

type accountsService struct {
	repo     repository.AccountsRepository
	txRepo   repository.TransactionsRepository
	producer *kafka.Producer
}

func NewAccountsService(
	r repository.AccountsRepository,
	txRepo repository.TransactionsRepository,
	prod *kafka.Producer,
) AccountsService {
	return &accountsService{
		repo:     r,
		txRepo:   txRepo,
		producer: prod,
	}
}

func (s *accountsService) GetAll(ctx context.Context, userID uint) ([]*entities.Account, error) {
	accounts, err := s.repo.GetAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %v", err)
	}

	return accounts, nil
}

func (s *accountsService) GetByID(ctx context.Context, userID, accountID uint) (*entities.Account, error) {
	account, err := s.repo.GetByID(ctx, userID, accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return account, nil
}

func (s *accountsService) Create(ctx context.Context, userID uint, req *entities.CreateAccountRequest) (*entities.Account, error) {
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

	if err := s.sendKafkaEvent(account); err != nil {
		lib.Log.Error("Failed to send Kafka event", zap.Error(err))
		return nil, err
	}

	return account, nil
}

func (s *accountsService) Deposit(ctx context.Context, userID, accountID uint, amount float64) (*entities.Account, error) {
	account, err := s.repo.GetByID(ctx, userID, accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("account not found")
		}
	}

	account.Balance += amount

	if err := s.repo.Update(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	tx := &entities.Transaction{
		FromAccountID: 0, // Внешний источник (например, банк)
		ToAccountID:   accountID,
		UserID:        userID,
		Amount:        amount,
		Description:   "Пополнение счёта",
		Type:          entities.Deposit,
		CreatedAt:     time.Now(),
	}

	if err := s.txRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *accountsService) Delete(ctx context.Context, userID, accountID uint) error {
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

// sendKafkaEvent отправляет событие в Kafka
func (s *accountsService) sendKafkaEvent(account *entities.Account) error {
	eventData := fmt.Sprintf(`{"account_id": %d, "user_id": %d, "status": "open"}`, account.ID, account.UserID)
	key := []byte(fmt.Sprint(account.ID))
	value := []byte(eventData)

	return s.producer.SendEvent(key, value)
}
