package repository

import (
	"bank-app-backend/internal/entities"
	"context"
	"gorm.io/gorm"
)

type AccountsRepository interface {
	GetAll(ctx context.Context, userID uint) ([]*entities.Account, error)
	GetByID(ctx context.Context, userID, accountID uint) (*entities.Account, error)
	Create(ctx context.Context, account *entities.Account) error
	Update(ctx context.Context, account *entities.Account) error
}

type accountsRepository struct {
	db *gorm.DB
}

func NewAccountsRepository(db *gorm.DB) AccountsRepository {
	return &accountsRepository{db: db}
}

func (r accountsRepository) GetAll(ctx context.Context, userID uint) ([]*entities.Account, error) {
	var accounts []*entities.Account

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r accountsRepository) GetByID(ctx context.Context, userID, accountID uint) (*entities.Account, error) {
	var account entities.Account

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND id = ?", userID, accountID).
		First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (r accountsRepository) Create(ctx context.Context, account *entities.Account) error {
	if err := r.db.WithContext(ctx).Create(&account).Error; err != nil {
		return err
	}

	return nil
}

func (r accountsRepository) Update(ctx context.Context, account *entities.Account) error {
	return r.db.WithContext(ctx).Save(account).Error
}
