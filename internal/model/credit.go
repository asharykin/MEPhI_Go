package model

import (
	"time"
)

const (
	CreditStatusActive   = "ACTIVE"
	CreditStatusClosed   = "CLOSED"
	CreditStatusOverdue  = "OVERDUE"
	CreditStatusRejected = "REJECTED"
)

type Credit struct {
	ID             int64     `json:"id" db:"id"`
	UserID         int64     `json:"user_id" db:"user_id"`
	AccountID      int64     `json:"account_id" db:"account_id"`
	Amount         float64   `json:"amount" db:"amount"`
	InterestRate   float64   `json:"interest_rate" db:"interest_rate"`
	TermMonths     int       `json:"term_months" db:"term_months"`
	MonthlyPayment float64   `json:"monthly_payment" db:"monthly_payment"`
	Status         string    `json:"status" db:"status"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type CreditCreateRequest struct {
	AccountID    int64   `json:"account_id" validate:"required"`
	Amount       float64 `json:"amount" validate:"required,gt=0"`
	InterestRate float64 `json:"interest_rate" validate:"required,gt=0"`
	TermMonths   int     `json:"term_months" validate:"required,gt=0"`
}

func (c *CreditCreateRequest) Validate() error {
	if c.AccountID <= 0 {
		return ErrInvalidAccountID
	}
	if !ValidateAmount(c.Amount) {
		return ErrInvalidAmount
	}
	if c.InterestRate <= 0 {
		return ErrInvalidInterestRate
	}
	if c.TermMonths <= 0 {
		return ErrInvalidTerm
	}
	return nil
}

type CreditResponse struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	AccountID      int64     `json:"account_id"`
	Amount         float64   `json:"amount"`
	InterestRate   float64   `json:"interest_rate"`
	TermMonths     int       `json:"term_months"`
	MonthlyPayment float64   `json:"monthly_payment"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

func (c *Credit) ToResponse() *CreditResponse {
	return &CreditResponse{
		ID:             c.ID,
		UserID:         c.UserID,
		AccountID:      c.AccountID,
		Amount:         c.Amount,
		InterestRate:   c.InterestRate,
		TermMonths:     c.TermMonths,
		MonthlyPayment: c.MonthlyPayment,
		Status:         c.Status,
		CreatedAt:      c.CreatedAt,
	}
}
