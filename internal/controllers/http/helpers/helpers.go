package helpers

import (
	"bank-app-backend/internal/entities"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

// ExtractUserID извлекает userID из контекста запроса
func ExtractUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, fmt.Errorf("unauthorized")
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		return 0, fmt.Errorf("invalid user ID")
	}

	return userIDUint, nil
}

// BuildTransactionFilter подтготавливает фильтр для транзацкии
func BuildTransactionFilter(c *gin.Context, userID uint) *entities.TransactionFilter {
	filter := &entities.TransactionFilter{
		UserID: userID,
		Page:   1,
		Limit:  20,
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filter.Page = page
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if fromStr := c.Query("fromDate"); fromStr != "" {
		if fromTime, err := time.Parse("2006-01-02", fromStr); err == nil {
			filter.FromDate = &fromTime
		}
	}

	if toStr := c.Query("toDate"); toStr != "" {
		if toTime, err := time.Parse("2006-01-02", toStr); err == nil {
			filter.ToDate = &toTime
		}
	}

	if t := c.Query("type"); t != "" {
		filter.Type = &t
	}

	if minAmount := c.Query("minAmount"); minAmount != "" {
		if v, err := strconv.ParseFloat(minAmount, 64); err == nil {
			filter.MinAmount = &v
		}
	}

	if maxAmount := c.Query("maxAmount"); maxAmount != "" {
		if v, err := strconv.ParseFloat(maxAmount, 64); err == nil {
			filter.MaxAmount = &v
		}
	}

	return filter
}
