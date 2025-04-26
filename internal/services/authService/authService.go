package authService

import (
	"bank-app-backend/internal/entities"
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
	Me(ctx context.Context, userID uint) (*entities.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (*entities.AuthResponse, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(r repository.AuthRepository) AuthService {
	return &authService{repo: r}
}

func (a *authService) RegisterUser(ctx context.Context, req entities.RegisterRequest) (*entities.User, error) {
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
	if err := a.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("could not create user: %v", err)
	}

	return user, nil
}

func (a *authService) Login(ctx context.Context, req entities.LoginRequest) (*entities.AuthResponse, error) {
	user, err := a.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("could not find user: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid pasword: %v", err)
	}

	accessToken, refreshToken, err := GenerateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("could not generate tokens: %v", err)
	}

	err = a.repo.SaveRefreshToken(ctx, user.ID, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("could not save refresh token: %v", err)
	}

	return &entities.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *authService) Logout(ctx context.Context, userID uint) error {
	return a.repo.DeleteRefreshToken(ctx, userID)
}

func (a *authService) Me(ctx context.Context, userID uint) (*entities.User, error) {
	user, err := a.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	return user, nil
}

func (a *authService) RefreshToken(ctx context.Context, refreshToken string) (*entities.AuthResponse, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return refreshSecret, nil
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

	user, err := a.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	accessToken, newRefreshToken, err := GenerateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("could not generate tokens: %v", err)
	}

	if err := a.repo.SaveRefreshToken(ctx, user.ID, newRefreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %v", err)
	}

	return &entities.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
