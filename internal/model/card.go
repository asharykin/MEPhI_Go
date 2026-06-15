package model

import (
	"banksystem/internal/crypto"
	"fmt"
	"time"
)

type Card struct {
	ID            int64     `json:"id"`
	AccountID     int64     `json:"account_id"`
	CardNumber    string    `json:"-"`
	ExpiryDate    time.Time `json:"-"`
	EncryptedData string    `json:"encrypted_data"`
	HashedCVV     string    `json:"hashed_cvv"`
	HMAC          string    `json:"hmac"`
	CreatedAt     time.Time `json:"created_at"`
}

type CardCreateRequest struct {
	AccountID  int64  `json:"account_id" validate:"required"`
	CardNumber string `json:"card_number" validate:"required"`
	ExpiryDate string `json:"expiry_date" validate:"required"`
	CVV        string `json:"cvv" validate:"required"`
}

type CardResponse struct {
	ID         int64     `json:"id"`
	AccountID  int64     `json:"account_id"`
	CardNumber string    `json:"card_number"` // Расшифрованный номер карты
	ExpiryDate string    `json:"expiry_date"` // Расшифрованная дата
	CreatedAt  time.Time `json:"created_at"`
}

func (c *CardCreateRequest) Validate() error {
	if c.AccountID <= 0 {
		return ErrInvalidAccountID
	}
	if !ValidateCardNumber(c.CardNumber) {
		return ErrInvalidCardNumber
	}
	if !ValidateExpiryDate(c.ExpiryDate) {
		return ErrInvalidExpiryDate
	}
	if !ValidateCVV(c.CVV) {
		return ErrInvalidCVV
	}
	return nil
}

func (c *Card) ToResponse(cardNumber, expiryDate string) CardResponse {
	return CardResponse{
		ID:         c.ID,
		AccountID:  c.AccountID,
		CardNumber: cardNumber,
		ExpiryDate: expiryDate,
		CreatedAt:  c.CreatedAt,
	}
}

func (c *Card) EncryptCardData() error {
	encrypted, err := crypto.EncryptCardData(c.CardNumber, c.ExpiryDate.Format("01/06"))
	if err != nil {
		return err
	}
	c.EncryptedData = encrypted
	c.HMAC = crypto.ComputeHMAC(encrypted)
	return nil
}

func (c *Card) DecryptCardData() error {
	if !crypto.VerifyHMAC(c.EncryptedData, c.HMAC) {
		return fmt.Errorf("неверный HMAC")
	}

	cardNumber, expiryDate, err := crypto.DecryptCardData(c.EncryptedData)
	if err != nil {
		return err
	}

	c.CardNumber = cardNumber
	expiry, err := time.Parse("01/06", expiryDate)
	if err != nil {
		return err
	}
	c.ExpiryDate = expiry

	return nil
}

func (c *Card) SetCVV(cvv string) error {
	hashed, err := crypto.HashCVV(cvv)
	if err != nil {
		return err
	}
	c.HashedCVV = hashed
	return nil
}

func (c *Card) VerifyCVV(cvv string) bool {
	return crypto.VerifyCVV(c.HashedCVV, cvv)
}
