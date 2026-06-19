package model

import (
	"time"
)

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
