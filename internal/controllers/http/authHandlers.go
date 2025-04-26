package http

import (
	"bank-app-backend/internal/entities"
	"bank-app-backend/internal/services/authService"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-jwt/jwt/v5"
	"net/http"
	_ "time"
)

type AuthHandler struct {
	service authService.AuthService
}

func NewAuthHandler(s authService.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// @Summary      Registering user
// @Description
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request body entities.RegisterRequest true "Register Request"
// @Success      201 {object} entities.User
// @Failure      400 {object} entities.ErrorResponse "Invalid input data"
// @Failure      500 {object} entities.ErrorResponse "Internal server error"
// @Router       /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req entities.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
	}

	user, err := h.service.RegisterUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

// @Summary      Login user
// @Description  Log in using email and password
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request body entities.LoginRequest true "User login request"
// @Success      200 {object} entities.AuthResponse "Login successful"
// @Failure      400 {object} entities.ErrorResponse "Invalid input data"
// @Failure      500 {object} entities.ErrorResponse "Internal server error"
// @Router       /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req entities.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
	}

	resp, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary      Refresh JWT token
// @Description  Refreshes the JWT access and refresh tokens
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request body entities.RefreshTokenRequest true "Refresh token request"
// @Success      200 {object} entities.AuthResponse "Tokens successfully refreshed"
// @Failure      400 {object} entities.ErrorResponse "Invalid refresh token"
// @Failure      500 {object} entities.ErrorResponse "Internal server error"
// @Router       /refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req entities.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	tokens, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// @Summary      Log out user
// @Description  Logs out the user by invalidating the session or token
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request body entities.LogoutRequest true "Logout request"
// @Success      200 {object} map[string]interface{} "Logout successful"
// @Failure      400 {object} entities.ErrorResponse "Invalid input data"
// @Failure      500 {object} entities.ErrorResponse "Internal server error"
// @Router       /logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req entities.LogoutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	err := h.service.Logout(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// @Summary      Get user information
// @Description  Retrieves the authenticated user's data
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200 {object} entities.User "User data retrieved successfully"
// @Failure      401 {object} entities.ErrorResponse "Unauthorized"
// @Failure      404 {object} entities.ErrorResponse "User not found"
// @Router       /me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := userIDInterface.(uint)

	user, err := h.service.Me(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	})
}
