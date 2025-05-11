package repository

import (
	"bank-app-backend/internal/entities"
	"context"
	"gorm.io/gorm"
)

type TransactionsRepository interface {
	Create(ctx context.Context, tx *entities.Transaction) error
	FindAll(ctx context.Context, filter *entities.TransactionFilter) ([]entities.Transaction, error)
	FindByID(ctx context.Context, id uint) (*entities.Transaction, error)
}

type transactionsRepository struct {
	db *gorm.DB
}

func NewTransactionsRepository(db *gorm.DB) TransactionsRepository {
	return &transactionsRepository{
		db: db,
	}
}

func (r *transactionsRepository) Create(ctx context.Context, tx *entities.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *transactionsRepository) FindAll(ctx context.Context, filter *entities.TransactionFilter) ([]entities.Transaction, error) {
	var txs []entities.Transaction
	db := r.db.WithContext(ctx).Model(&entities.Transaction{}).Where("user_id = ?", filter.UserID)

	if filter.Type != nil {
		db = db.Where("type = ?", *filter.Type)
	}
	if filter.FromDate != nil {
		db = db.Where("created_at >= ?", *filter.FromDate)
	}
	if filter.ToDate != nil {
		db = db.Where("created_at <= ?", *filter.ToDate)
	}
	if filter.MinAmount != nil {
		db = db.Where("amount >= ?", *filter.MinAmount)
	}
	if filter.MaxAmount != nil {
		db = db.Where("amount <= ?", *filter.MaxAmount)
	}

	limit := filter.Limit
	offset := (filter.Page - 1) * filter.Limit

	err := db.Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&txs).Error

	return txs, err
}

func (r *transactionsRepository) FindByID(ctx context.Context, id uint) (*entities.Transaction, error) {
	var tx entities.Transaction
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}
