package http

import (
	"bank-app-backend/internal/controllers/http/helpers"
	"bank-app-backend/internal/entities"
	"bank-app-backend/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type UsersHandler struct {
	service services.UsersService
}

func NewUsersHandler(s services.UsersService) *UsersHandler {
	return &UsersHandler{service: s}
}

// @Summary      Get user information
// @Description  Retrieves the authenticated user's data
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200 {object} entities.UserResponse "User data retrieved successfully"
// @Failure      401 {object} entities.ErrorResponse "Unauthorized"
// @Failure      404 {object} entities.ErrorResponse "User not found"
// @Router       /me [get]
func (h *UsersHandler) Me(c *gin.Context) {
	userID, err := helpers.ExtractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

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

// @Summary      Get all users
// @Description  Retrieves all users
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200 {array} entities.UserResponse "Users retrieved successfully"
// @Failure      500 {object} entities.ErrorResponse "Internal server error"
// @Router       /users [get]
func (h *UsersHandler) GetAll(c *gin.Context) {
	users, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	for _, user := range users {
		user.Password = ""
	}

	c.JSON(http.StatusOK, entities.UsersToResponse(users))
}

// @Summary      Update a user
// @Description  Updates a user's information by ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      int                     true  "User ID"
// @Param        user  body      entities.UpdateUserRequest true  "Updated user information"
// @Success      200   {object}  entities.UserResponse   "User updated successfully"
// @Failure      400   {object}  entities.ErrorResponse   "Invalid input or user ID"
// @Failure      404   {object}  entities.ErrorResponse   "User not found"
// @Failure      500   {object}  entities.ErrorResponse   "Internal server error"
// @Router       /users/{id} [patch]
func (h *UsersHandler) Update(c *gin.Context) {
	var input entities.UpdateUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.service.Update(c.Request.Context(), uint(userID), &input)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}
