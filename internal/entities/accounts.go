package entities

import "time"

// Account represents the database model for a user's account.
// @Description Account entity containing balance, currency, and status information.
// @example { "id": 1, "user_id": 2, "type": "deposit", "currency": "RUB", "balance": 1000.50, "status": "active" }
type Account struct {
	ID        uint    `gorm:"primary_key;auto_increment"`
	UserID    uint    `gorm:"primary_key;not null"`
	Type      string  `gorm:"not null"`
	Currency  string  `gorm:"not null"`
	Balance   float64 `gorm:"not null"`
	Status    string  `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateAccountRequest represents the payload required to create a new account.
// @Description Request payload for creating a new user account with a specified type and currency.
// @example { "type": "deposit", "currency": "RUB" }
type CreateAccountRequest struct {
	Type     string `json:"type" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

// UpdateAccount represents the payload to update fields of an existing account.
// @Description Request payload to update account type, currency, or status.
// @example { "type": "savings", "currency": "USD", "status": "closed" }
type UpdateAccount struct {
	Type     *string `json:"type"`
	Currency *string `json:"currency"`
	Status   *string `json:"status"`
}

// AccountResponse represents the public response structure of an account.
// @Description Response returned when retrieving account information.
// @example { "id": 1, "user_id": 2, "type": "deposit", "currency": "RUB", "balance": 1000.50, "status": "active" }
type AccountResponse struct {
	ID       uint    `json:"id"`
	UserID   uint    `json:"user_id"`
	Type     string  `json:"type"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
	Status   string  `json:"status"`
}

// DepositRequest представляет тело запроса для пополнения счёта.
// @Description Запрос для пополнения счёта пользователя на определённую сумму.
// @example { "account_id": 1, "amount": 1000.50 }
type DepositRequest struct {
	AccountID uint    `json:"account_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
}

// MessageResponse represents a success message response.
// @Description Success message response
// @example { "message": "Account closed successfully" }
type MessageResponse struct {
	Message string `json:"message"`
}

func (a *Account) ToResponse() *AccountResponse {
	return &AccountResponse{
		ID:       a.ID,
		UserID:   a.UserID,
		Type:     a.Type,
		Currency: a.Currency,
		Balance:  a.Balance,
		Status:   a.Status,
	}
}

func AccountsToResponse(accounts []*Account) []*AccountResponse {
	responses := make([]*AccountResponse, len(accounts))
	for i, account := range accounts {
		responses[i] = account.ToResponse()
	}
	return responses
}
