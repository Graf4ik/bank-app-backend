package repository

import (
	"bank-app-backend/internal/entities"
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByID(ctx context.Context, id uint) (*entities.User, error)
	SaveRefreshToken(ctx context.Context, userID uint, token string) error
	FindUserByRefreshToken(ctx context.Context, token string) (*entities.User, error)
	DeleteRefreshToken(ctx context.Context, userID uint) error
	DeleteExpiredTokens(ctx context.Context) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (a *authRepository) CreateUser(ctx context.Context, user *entities.User) error {
	return a.db.WithContext(ctx).Create(user).Error
}

func (a *authRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	if err := a.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *authRepository) FindByID(ctx context.Context, id uint) (*entities.User, error) {
	var user entities.User
	if err := a.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *authRepository) SaveRefreshToken(ctx context.Context, userID uint, token string) error {
	refreshToken := &entities.RefreshToken{
		Token:  token,
		UserID: fmt.Sprint(userID),
	}

	return a.db.WithContext(ctx).Create(refreshToken).Error
}

func (a *authRepository) FindUserByRefreshToken(ctx context.Context, token string) (*entities.User, error) {
	var user entities.User
	if err := a.db.WithContext(ctx).Where("refresh_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *authRepository) DeleteRefreshToken(ctx context.Context, userID uint) error {
	return a.db.WithContext(ctx).
		Where("id = ?", userID).
		Delete(&entities.RefreshToken{}).Error
}

func (a *authRepository) DeleteExpiredTokens(ctx context.Context) error {
	return a.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&entities.RefreshToken{}).Error
}
