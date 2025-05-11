package entities

import "time"

type TransferType string

const (
	InternalTransfer TransferType = "internal"
	ExternalTransfer TransferType = "external"
	Deposit          TransferType = "deposit"
)

// TransferRequest represents a request to initiate a transfer between accounts.
// @Description TransferRequest is used to initiate a transfer between two accounts.
// @Model
type TransferRequest struct {
	UserID        uint         `json:"user_id"`
	FromAccountID uint         `json:"from_account_id"`
	ToAccountID   uint         `json:"to_account_id"`
	Amount        float64      `json:"amount"`
	Description   string       `json:"description,omitempty"`
	Type          TransferType `json:"type"`
}

// Transaction represents a financial transaction.
// @Description Transaction is a record of a financial transaction between accounts.
// @Model
type Transaction struct {
	ID            uint         `json:"id"`
	UserID        uint         `json:"user_id"`
	FromAccountID uint         `json:"from_account_id"`
	ToAccountID   uint         `json:"to_account_id"`
	Amount        float64      `json:"amount"`
	Description   string       `json:"description"`
	Type          TransferType `json:"type"`
	CreatedAt     time.Time    `json:"created_at"`
}

// TransactionFilter is used to filter transactions by different parameters.
// @Description TransactionFilter is used to filter transactions based on criteria like date, amount, and type.
// @Model
type TransactionFilter struct {
	UserID    uint       `json:"user_id"`
	FromDate  *time.Time `json:"from_date"`
	ToDate    *time.Time `json:"to_date"`
	Type      *string    `json:"type"`
	MinAmount *float64   `json:"min_amount"`
	MaxAmount *float64   `json:"max_amount"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
}
