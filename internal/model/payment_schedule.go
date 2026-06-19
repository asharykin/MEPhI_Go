package model

import (
	"time"
)

type PaymentSchedule struct {
	ID          int       `json:"id"`
	CreditID    int       `json:"credit_id"`
	PaymentDate time.Time `json:"payment_date"`
	Amount      float64   `json:"amount"`
	Principal   float64   `json:"principal"`
	Interest    float64   `json:"interest"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PaymentScheduleResponse struct {
	ID          int       `json:"id"`
	CreditID    int       `json:"credit_id"`
	PaymentDate time.Time `json:"payment_date"`
	Amount      float64   `json:"amount"`
	Principal   float64   `json:"principal"`
	Interest    float64   `json:"interest"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (p *PaymentSchedule) ToResponse() *PaymentScheduleResponse {
	return &PaymentScheduleResponse{
		ID:          p.ID,
		CreditID:    p.CreditID,
		PaymentDate: p.PaymentDate,
		Amount:      p.Amount,
		Principal:   p.Principal,
		Interest:    p.Interest,
		Status:      p.Status,
		CreatedAt:   p.CreatedAt,
	}
}
