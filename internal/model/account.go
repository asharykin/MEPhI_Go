package model

import (
	"time"
)

type Account struct {
	ID        string
	UserID    string
	Balance   float64
	Currency  string
	CreatedAt time.Time
}
