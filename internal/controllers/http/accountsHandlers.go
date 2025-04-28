package http

import (
	"bank-app-backend/internal/entities"
	"bank-app-backend/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AccountsHandler struct {
	service services.AccountsService
}

func NewAccountsHandler(s services.AccountsService) *AccountsHandler {
	return &AccountsHandler{service: s}
}

// GetAllByUser godoc
// @Summary Get all user accounts
// @Description Returns a list of accounts belonging to the authenticated user
// @Tags accounts
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entities.AccountResponse
// @Failure 401 {object} entities.ErrorResponse "Unauthorized"
// @Failure 500 {object} entities.ErrorResponse "Internal server error"
// @Router /auth/accounts [get]
func (h *AccountsHandler) GetAllByUser(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	userID := userIDInterface.(uint)

	accounts, err := h.service.GetAll(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, entities.AccountsToResponse(accounts))
}

// GetByID godoc
// @Summary Get a user account by ID
// @Description Returns a specific account by ID if it belongs to the authenticated user
// @Tags accounts
// @Security BearerAuth
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} entities.AccountResponse
// @Failure 400 {object} entities.ErrorResponse "Invalid account ID"
// @Failure 401 {object} entities.ErrorResponse "Unauthorized"
// @Failure 404 {object} entities.ErrorResponse "Account not found"
// @Router /auth/accounts/{id} [get]
func (h *AccountsHandler) GetByID(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := userIDInterface.(uint)

	accountIDParam := c.Param("id")
	accountIDUint64, err := strconv.ParseUint(accountIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	accountID := uint(accountIDUint64)

	account, err := h.service.GetByID(c.Request.Context(), userID, accountID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
	}

	c.JSON(http.StatusOK, account.ToResponse())
}

// Create godoc
// @Summary Create a new account
// @Description Creates a new account for the authenticated user
// @Tags accounts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param account body entities.CreateAccountRequest true "Account creation data"
// @Success 201 {object} entities.AccountResponse
// @Failure 400 {object} entities.ErrorResponse "Invalid input"
// @Failure 401 {object} entities.ErrorResponse "Unauthorized"
// @Failure 500 {object} entities.ErrorResponse "Failed to create account"
// @Router /auth/accounts [post]
func (h *AccountsHandler) Create(c *gin.Context) {
	var req entities.CreateAccountRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
	userID := userIDInterface.(uint)

	account, err := h.service.Create(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account", "details": err.Error()})
	}

	c.JSON(http.StatusCreated, account.ToResponse())
}

// CloseAccount godoc
// @Summary Close a user account
// @Description Closes an account if its balance is zero
// @Tags accounts
// @Security BearerAuth
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} entities.MessageResponse "Account closed successfully"
// @Failure 400 {object} entities.ErrorResponse "Cannot close account"
// @Failure 401 {object} entities.ErrorResponse "Unauthorized"
// @Router /auth/accounts/{id} [patch]
func (h *AccountsHandler) CloseAccount(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := userIDInterface.(uint)

	accountIDParam := c.Param("id")
	accountIDUint64, err := strconv.ParseUint(accountIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	accountID := uint(accountIDUint64)

	err = h.service.Delete(c.Request.Context(), userID, accountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account closed successfully"})
}
