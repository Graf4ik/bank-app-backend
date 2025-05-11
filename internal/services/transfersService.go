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
	"time"
)

type TransfersService interface {
	ProcessTransfer(ctx context.Context, req entities.TransferRequest) (*entities.Transaction, error)
}

type transfersService struct {
	txRepo   repository.TransactionsRepository
	accRepo  repository.AccountsRepository
	producer *kafka.Producer
}

func NewTransfersService(txRepo repository.TransactionsRepository, accRepo repository.AccountsRepository, prod *kafka.Producer) TransfersService {
	return &transfersService{
		txRepo:   txRepo,
		accRepo:  accRepo,
		producer: prod,
	}
}

func (s *transfersService) ProcessTransfer(ctx context.Context, req entities.TransferRequest) (*entities.Transaction, error) {
	fromAccount, err := s.accRepo.GetByID(ctx, req.UserID, req.FromAccountID)
	if err != nil {
		return nil, err
	}

	if fromAccount.Balance < req.Amount {
		return nil, errors.New("insufficient funds")
	}

	toAccount, err := s.accRepo.GetByID(ctx, req.UserID, req.ToAccountID)
	if err != nil {
		return nil, err
	}
	fmt.Printf("REQ: %+v\n", req)
	fromAccount.Balance -= req.Amount
	toAccount.Balance += req.Amount

	if err := s.accRepo.Update(ctx, fromAccount); err != nil {
		return nil, err
	}
	if err := s.accRepo.Update(ctx, toAccount); err != nil {
		return nil, err
	}

	tx := &entities.Transaction{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		UserID:        req.UserID,
		Amount:        req.Amount,
		Description:   req.Description,
		Type:          req.Type,
		CreatedAt:     time.Now(),
	}

	if err := s.sendKafkaEvent(tx); err != nil {
		lib.Log.Error("Failed to send Kafka event", zap.Error(err))
		return nil, err
	}

	if err := s.txRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// sendKafkaEvent отправляет событие в Kafka
func (s *transfersService) sendKafkaEvent(tx *entities.Transaction) error {
	txEvent := fmt.Sprintf("Transaction of %f completed for user %d", tx.Amount, tx.UserID)
	return s.producer.SendEvent([]byte(fmt.Sprintf("%d", tx.ID)), []byte(txEvent))
}
