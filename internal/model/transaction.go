package model

import (
	"time"
)

type Transaction struct {
	ID          string
	SenderID    string
	ReceiverID  string
	Amount      float64
	Type        string
	Description string
	CreatedAt   time.Time
}
