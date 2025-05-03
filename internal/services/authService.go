package services

import (
	"bank-app-backend/internal/entities"
	redis "bank-app-backend/internal/lib/redis"
	"bank-app-backend/internal/lib/token"
	"bank-app-backend/internal/repository"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(ctx context.Context, req entities.RegisterRequest) (*entities.User, error)
	Login(ctx context.Context, req entities.LoginRequest) (*entities.AuthResponse, error)
	Logout(ctx context.Context, userID uint) error
	RefreshToken(ctx context.Context, refreshToken string) (*entities.AuthResponse, error)
}

type authService struct {
	repo  repository.UsersRepository
	redis *redis.Client
}

func NewAuthService(r repository.UsersRepository, redisClient *redis.Client) AuthService {
	return &authService{
		repo:  r,
		redis: redisClient,
	}
}

func (s *authService) RegisterUser(ctx context.Context, req entities.RegisterRequest) (*entities.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %v", err)
	}

	user := &entities.User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
	}

	fmt.Println("Password from DB:", user.Password)
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("could not create user: %v", err)
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, req entities.LoginRequest) (*entities.AuthResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid pasword: %v", err)
	}

	accessToken, refreshToken, err := lib.GenerateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("could not generate tokens: %v", err)
	}

	err = s.repo.SaveRefreshToken(ctx, user.ID, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("could not save refresh token: %v", err)
	}

	return &entities.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID uint) error {
	return s.repo.DeleteRefreshToken(ctx, userID)
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*entities.AuthResponse, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return lib.RefreshSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] == nil {
		return nil, fmt.Errorf("invalid token claims")
	}

	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}
	userID := uint(userIDFloat)

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	accessToken, newRefreshToken, err := lib.GenerateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("could not generate tokens: %v", err)
	}

	if err := s.repo.SaveRefreshToken(ctx, user.ID, newRefreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %v", err)
	}

	return &entities.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
