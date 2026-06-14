package models

import "time"

type Payment struct {
	ID        int64     `json:"id"`
	CreditID  int64     `json:"credit_id"`
	Amount    float64   `json:"amount"`
	DueDate   time.Time `json:"due_date"`
	Status    string    `json:"status"` // pending, paid, overdue
	CreatedAt time.Time `json:"created_at"`
}

type PaymentCreateRequest struct {
	CreditID int64     `json:"credit_id" validate:"required"`
	Amount   float64   `json:"amount" validate:"required,gt=0"`
	DueDate  time.Time `json:"due_date" validate:"required"`
}

type PaymentResponse struct {
	ID        int64     `json:"id"`
	CreditID  int64     `json:"credit_id"`
	Amount    float64   `json:"amount"`
	DueDate   time.Time `json:"due_date"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *Payment) ToResponse() PaymentResponse {
	return PaymentResponse{
		ID:        p.ID,
		CreditID:  p.CreditID,
		Amount:    p.Amount,
		DueDate:   p.DueDate,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
	}
}
