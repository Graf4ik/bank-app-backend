package repository

import (
	"bank-app-backend/internal/entities"
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByID(ctx context.Context, id uint) (*entities.User, error)
	FindAll(ctx context.Context) ([]*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	SaveRefreshToken(ctx context.Context, userID uint, token string) error
	FindUserByRefreshToken(ctx context.Context, token string) (*entities.User, error)
	DeleteRefreshToken(ctx context.Context, userID uint) error
	DeleteExpiredTokens(ctx context.Context) error
}

type usersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) UsersRepository {
	return &usersRepository{db: db}
}

func (r *usersRepository) CreateUser(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *usersRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *usersRepository) FindByID(ctx context.Context, id uint) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *usersRepository) FindAll(ctx context.Context) ([]*entities.User, error) {
	var users []*entities.User

	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *usersRepository) Update(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *usersRepository) SaveRefreshToken(ctx context.Context, userID uint, token string) error {
	refreshToken := &entities.RefreshToken{
		Token:  token,
		UserID: fmt.Sprint(userID),
	}

	return r.db.WithContext(ctx).Create(refreshToken).Error
}

func (r *usersRepository) FindUserByRefreshToken(ctx context.Context, token string) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).Where("refresh_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *usersRepository) DeleteRefreshToken(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ?", userID).
		Delete(&entities.RefreshToken{}).Error
}

func (r *usersRepository) DeleteExpiredTokens(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&entities.RefreshToken{}).Error
}
