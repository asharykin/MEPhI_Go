package models

import "time"

type CreditPayment struct {
	ID        int64     `json:"id"`
	CreditID  int64     `json:"credit_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	DueDate   time.Time `json:"due_date"`
	CreatedAt time.Time `json:"created_at"`
}

type CreditPaymentCreateRequest struct {
	CreditID int64     `json:"credit_id"`
	Amount   float64   `json:"amount"`
	DueDate  time.Time `json:"due_date"`
}

type CreditPaymentResponse struct {
	ID        int64     `json:"id"`
	CreditID  int64     `json:"credit_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	DueDate   time.Time `json:"due_date"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *CreditPayment) ToResponse() *CreditPaymentResponse {
	return &CreditPaymentResponse{
		ID:        p.ID,
		CreditID:  p.CreditID,
		Amount:    p.Amount,
		Status:    p.Status,
		DueDate:   p.DueDate,
		CreatedAt: p.CreatedAt,
	}
}
