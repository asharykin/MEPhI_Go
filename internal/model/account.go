package model

import (
	"database/sql"
	"time"
)

const (
	AccountTypeChecking = "checking"
	AccountTypeDebit    = "savings"
	AccountTypeCredit   = "credit"
)

type Account struct {
	ID        int64        `json:"id"`
	UserID    int64        `json:"user_id"`
	Balance   float64      `json:"balance"`
	Type      string       `json:"type"` // checking, savings, credit
	IsActive  bool         `json:"is_active"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

type AccountCreateRequest struct {
	UserID int64  `json:"user_id" validate:"required"`
	Type   string `json:"type" validate:"required,oneof=checking savings credit"`
}

func (a *AccountCreateRequest) Validate() error {
	if a.UserID <= 0 {
		return ErrInvalidUserID
	}
	if a.Type != "checking" && a.Type != "savings" && a.Type != "credit" {
		return ErrInvalidAccountType
	}
	return nil
}

type AccountTransferRequest struct {
	FromAccountID int     `json:"from_account_id" validate:"required"`
	ToAccountID   int     `json:"to_account_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
}

func (a *AccountTransferRequest) Validate() error {
	if a.FromAccountID <= 0 || a.ToAccountID <= 0 {
		return ErrInvalidAccountID
	}
	if !ValidateAmount(a.Amount) {
		return ErrInvalidAmount
	}
	return nil
}

type AccountResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Balance   float64   `json:"balance"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *Account) ToResponse() *AccountResponse {
	return &AccountResponse{
		ID:        a.ID,
		UserID:    a.UserID,
		Balance:   a.Balance,
		Type:      a.Type,
		CreatedAt: a.CreatedAt.Time,
	}
}
