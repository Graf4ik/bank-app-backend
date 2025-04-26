package entities

import "time"

// @Description RegisterRequest model
// @example { "email": "user@example.com", "username": "user1", "password": "123456" }
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Description LoginRequest model
// @example { "email": "user@example.com", "password": "123456" }
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// @Description RefreshTokenRequest model
// @example { "refresh_token": "old_refresh_token_value" }
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// @Description LogoutRequest model
// @example { "userId": 123 }
type LogoutRequest struct {
	UserID uint `json:"userId" binding:"required"`
}

// @Description RefreshTokenRequest model
// @example { "refresh_token": "old_refresh_token_value" }
type RefreshToken struct {
	ID        uint   `gorm:"primaryKey"`
	Token     string `gorm:"uniqueIndex;not null"`
	UserID    string `gorm:"not null"`
	CreatedAt time.Time
	ExpiresAt time.Time `gorm:"index"`
}

// @Description AuthResponse contains the access and refresh tokens
// @example { "accessToken": "new_access_token_value", "refreshToken": "new_refresh_token_value" }
type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// User represents a user in the system.
// @Description User model
// @example { "id": 1, "email": "user@example.com", "username": "user1", password: "123456" }
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	Username string
	Password string
}

// ErrorResponse represents the error that is returned when the request fails
// @Description ErrorResponse structure
// @example { "error": "Invalid request", "details": "Detailed error message" }
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}
