package model

import "time"

type Account struct {
	ID        string
	UserID    string
	Balance   float64
	Currency  string
	CreatedAt time.Time
}

type Card struct {
	ID              string
	UserID          string
	AccountID       string
	NumberEncrypted []byte
	ExpiryEncrypted []byte
	CVVHash         []byte
	HMAC            string
	CreatedAt       time.Time
}

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

type PaymentSchedule struct {
	ID             string
	CreditID       string
	PaymentDate    time.Time
	Amount         float64
	IsPaid         bool
	LateFeeApplied bool
}

type Transaction struct {
	ID          string
	SenderID    string
	ReceiverID  string
	Amount      float64
	Type        string
	Description string
	CreatedAt   time.Time
}

type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
