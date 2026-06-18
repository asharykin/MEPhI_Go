package model

import (
	"database/sql"
	"time"
)

const (
	TransactionTypeDeposit    = "DEPOSIT"
	TransactionTypeWithdrawal = "WITHDRAWAL"
	TransactionTypeTransfer   = "TRANSFER"
	TransactionTypePayment    = "PAYMENT"

	TransactionStatusPending   = "PENDING"
	TransactionStatusCompleted = "COMPLETED"
	TransactionStatusFailed    = "FAILED"
)

type Transaction struct {
	ID          int64         `json:"id"`
	AccountID   int64         `json:"account_id"`
	Type        string        `json:"type"` // deposit, withdraw, transfer_in, transfer_out
	Amount      float64       `json:"amount"`
	Status      string        `json:"status"`
	ToAccountID sql.NullInt64 `json:"to_account_id,omitempty"`
	CreatedAt   sql.NullTime  `json:"created_at"`
}

type TransactionCreateRequest struct {
	AccountID   int64   `json:"account_id" validate:"required"`
	Type        string  `json:"type" validate:"required,oneof=deposit withdraw transfer_in transfer_out"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	ToAccountID int64   `json:"to_account_id,omitempty"`
}

type TransactionResponse struct {
	ID          int64     `json:"id"`
	AccountID   int64     `json:"account_id"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	ToAccountID int64     `json:"to_account_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func (t *Transaction) ToResponse() TransactionResponse {
	var toAccountID int64
	if t.ToAccountID.Valid {
		toAccountID = t.ToAccountID.Int64
	}

	var createdAt time.Time
	if t.CreatedAt.Valid {
		createdAt = t.CreatedAt.Time
	}

	return TransactionResponse{
		ID:          t.ID,
		AccountID:   t.AccountID,
		Type:        t.Type,
		Amount:      t.Amount,
		Status:      t.Status,
		ToAccountID: toAccountID,
		CreatedAt:   createdAt,
	}
}

func ValidateTransactionType(tType string) bool {
	switch tType {
	case TransactionTypeDeposit,
		TransactionTypeWithdrawal,
		TransactionTypeTransfer,
		TransactionTypePayment:
		return true
	default:
		return false
	}
}

func ValidateTransactionStatus(status string) bool {
	switch status {
	case TransactionStatusPending,
		TransactionStatusCompleted,
		TransactionStatusFailed:
		return true
	default:
		return false
	}
}
