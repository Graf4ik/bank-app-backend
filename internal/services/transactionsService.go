package services

import (
	"bank-app-backend/internal/entities"
	"bank-app-backend/internal/repository"
	"context"
)

type TransactionsService interface {
	GetTransactions(ctx context.Context, filter *entities.TransactionFilter) ([]entities.Transaction, error)
	GetTransactionByID(ctx context.Context, id uint) (*entities.Transaction, error)
}

type transactionsService struct {
	txRepo repository.TransactionsRepository
}

func NewTransactionService(txRepo repository.TransactionsRepository) TransactionsService {
	return &transactionsService{txRepo: txRepo}
}

func (s *transactionsService) GetTransactions(ctx context.Context, filter *entities.TransactionFilter) ([]entities.Transaction, error) {
	return s.txRepo.FindAll(ctx, filter)
}

func (s *transactionsService) GetTransactionByID(ctx context.Context, id uint) (*entities.Transaction, error) {
	return s.txRepo.FindByID(ctx, id)
}
