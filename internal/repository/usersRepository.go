package repository

import (
	"bank-app-backend/internal/entities"
	lib "bank-app-backend/internal/lib/logger"
	redis "bank-app-backend/internal/lib/redis"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
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
	db    *gorm.DB
	redis *redis.Client
}

func NewUsersRepository(db *gorm.DB, redisClient *redis.Client) UsersRepository {
	return &usersRepository{
		db:    db,
		redis: redisClient,
	}
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
	key := fmt.Sprintf("user:%d", id)

	cached, err := r.redis.Get(ctx, key)
	if err == nil {
		var user entities.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			lib.Log.Info("User found in Redis", zap.Uint("user_id", id))
			return &user, nil
		}
	}

	var user entities.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		lib.Log.Error("User not found in DB", zap.Uint("user_id", id), zap.Error(err))
		return nil, err
	}

	data, err := json.Marshal(&user)
	if err == nil {
		lib.Log.Info("User cached in Redis", zap.Uint("user_id", id))
		err := r.redis.Set(ctx, key, data, 10*time.Minute)
		if err != nil {
			return nil, err
		}
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
	key := fmt.Sprintf("refresh_token:%d", userID)
	if err := r.redis.Set(ctx, key, token, 7*24*time.Hour); err != nil {
		lib.Log.Error("Failed to save refresh token", zap.Uint("userID", userID), zap.Error(err))
		return err
	}

	return nil
}

func (r *usersRepository) FindUserByRefreshToken(ctx context.Context, token string) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).Where("refresh_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *usersRepository) DeleteRefreshToken(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("refresh_token:%d", userID)
	err := r.redis.Del(ctx, key)
	if err != nil {
		lib.Log.Error("Error deleting key", zap.String("key", key), zap.Error(err))
	} else {
		lib.Log.Info("Key successfully deleted", zap.String("key", key))
	}
	return err
}

func (r *usersRepository) DeleteExpiredTokens(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&entities.RefreshToken{}).Error
}
