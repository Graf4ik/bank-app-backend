package http

import (
	"bank-app-backend/internal/controllers/http/helpers"
	"bank-app-backend/internal/entities"
	"bank-app-backend/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TransactionsHandler struct {
	transfersService services.TransfersService
	txService        services.TransactionsService
}

func NewTransactionsHandler(transferService services.TransfersService, txService services.TransactionsService) *TransactionsHandler {
	return &TransactionsHandler{
		transfersService: transferService,
		txService:        txService,
	}
}

// @Tags Transactions
// @Summary Internal transfer
// @Description Process an internal transfer between accounts
// @Accept json
// @Produce json
// @Param transfer body entities.TransferRequest true "Transfer request"
// @Success 200 {object} entities.Transaction "Transaction details"
// @Failure 400 {object} entities.ErrorResponse "Error processing transfer"
// @Failure 401 {object} entities.ErrorResponse "Unauthorized"
// @Router /auth/transfers/internal [post]
func (h *TransactionsHandler) InternalTransfer(c *gin.Context) {
	var req entities.TransferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := helpers.ExtractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	req.UserID = userID
	req.Type = entities.InternalTransfer
	tx, err := h.transfersService.ProcessTransfer(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tx})
}

// @Tags Transactions
// @Summary External transfer
// @Description Process an external transfer between accounts
// @Accept json
// @Produce json
// @Param transfer body entities.TransferRequest true "Transfer request"
// @Success 200 {object} entities.Transaction "Transaction details"
// @Failure 400 {object} entities.ErrorResponse "Error processing transfer"
// @Failure 401 {object} entities.ErrorResponse "Unauthorized"
// @Router /auth/transfers/external [post]
func (h *TransactionsHandler) ExternalTransfer(c *gin.Context) {
	var req entities.TransferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userID, err := helpers.ExtractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	req.UserID = userID
	req.Type = entities.ExternalTransfer
	tx, err := h.transfersService.ProcessTransfer(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}

	c.JSON(http.StatusOK, tx)
}

// @Tags Transactions
// @Summary Get a list of transactions
// @Description Get a list of transactions for a user, with optional filters for pagination, date range, type, and amount
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit number of transactions"
// @Param fromDate query string false "From date (YYYY-MM-DD)"
// @Param toDate query string false "To date (YYYY-MM-DD)"
// @Param type query string false "Transaction type"
// @Param minAmount query float64 false "Minimum amount"
// @Param maxAmount query float64 false "Maximum amount"
// @Success 200 {array} entities.Transaction "List of transactions"
// @Failure 401 {object} entities.ErrorResponse "Unauthorized"
// @Failure 400 {object} entities.ErrorResponse "Invalid request"
// @Router /auth/transactions [get]
func (h *TransactionsHandler) GetTransactions(c *gin.Context) {
	userID, err := helpers.ExtractUserID(c)
	filter := helpers.BuildTransactionFilter(c, userID)

	txs, err := h.txService.GetTransactions(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	c.JSON(http.StatusOK, txs)
}

// @Tags Transactions
// @Summary Get a transaction by ID
// @Description Get transaction details by ID
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} entities.Transaction "Transaction details"
// @Failure 404 {object} entities.ErrorResponse "Transaction not found"
// @Failure 400 {object} entities.ErrorResponse "Invalid transaction ID"
// @Router /auth/transactions/{id} [get]
func (h *TransactionsHandler) GetTransactionById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	tx, err := h.txService.GetTransactionByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
	}

	c.JSON(http.StatusOK, tx)
}
