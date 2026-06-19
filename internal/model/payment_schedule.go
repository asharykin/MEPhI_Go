package model

import (
	"time"
)

type PaymentSchedule struct {
	ID             string
	CreditID       string
	PaymentDate    time.Time
	Amount         float64
	IsPaid         bool
	LateFeeApplied bool
}
