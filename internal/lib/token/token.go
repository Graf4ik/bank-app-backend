package token

import (
	"bank-app-backend/internal/entities"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	// TODO: вынести в env
	accessSecret  = []byte("jwt_access_secret")
	refreshSecret = []byte("jwt_refresh_secret")
)

func GenerateTokens(user *entities.User) (accessToken string, refreshToken string, err error) {
	accessToken, err = generateJWT(user, accessSecret, time.Minute*15)
	if err != nil {
		return
	}
	refreshToken, err = generateJWT(user, refreshSecret, time.Hour*24*7)
	return
}

func generateJWT(user *entities.User, secret []byte, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"exp":   time.Now().Add(ttl).Unix(),
		"email": user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
