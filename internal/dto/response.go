package dto

import "time"

type AccountResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type AnalyticsResponse struct {
	Income           float64 `json:"income"`
	Expense          float64 `json:"expense"`
	CreditLoad       float64 `json:"credit_load"`
	PredictedBalance float64 `json:"predicted_balance"`
}

type AuthResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

type CardResponse struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	AccountID      string    `json:"account_id"`
	NumberLastFour string    `json:"number_last_four"`
	ExpiryDate     string    `json:"expiry_date"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreditResponse struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	AccountID      string    `json:"account_id"`
	Principal      float64   `json:"principal"`
	InterestRate   float64   `json:"interest_rate"`
	TermMonths     int       `json:"term_months"`
	MonthlyPayment float64   `json:"monthly_payment"`
	RemainingDebt  float64   `json:"remaining_debt"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

type PaymentScheduleResponse struct {
	ID             string    `json:"id"`
	CreditID       string    `json:"credit_id"`
	PaymentDate    time.Time `json:"payment_date"`
	Amount         float64   `json:"amount"`
	IsPaid         bool      `json:"is_paid"`
	LateFeeApplied bool      `json:"late_fee_applied"`
}

type PredictBalanceResponse struct {
	AccountID        string  `json:"account_id"`
	CurrentBalance   float64 `json:"current_balance"`
	PredictedBalance float64 `json:"predicted_balance"`
	Days             int     `json:"days"`
}

type ErrorResponse struct {
	Message string
}

func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
	}
}

func (r *ErrorResponse) Error() string {
	return r.Message
}
