package model

import (
	"time"
)

type Credit struct {
	ID             string
	UserID         string
	AccountID      string
	Principal      float64
	InterestRate   float64
	TermMonths     int
	MonthlyPayment float64
	RemainingDebt  float64
	Status         string
	CreatedAt      time.Time
}
