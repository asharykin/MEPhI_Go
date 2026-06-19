package dto

import "strings"

type CreateAccountRequest struct {
	Currency string `json:"currency"`
}

func (r *CreateAccountRequest) Validate() error {
	if r.Currency != "RUB" {
		return NewErrorResponse("Currency must be RUB")
	}

	return nil
}

type CreateCardRequest struct {
	AccountID string `json:"account_id"`
}

func (r *CreateCardRequest) Validate() error {
	if r.AccountID == "" {
		return NewErrorResponse("Account ID is required")
	}

	return nil
}

type CreateCreditRequest struct {
	AccountID  string  `json:"account_id"`
	Principal  float64 `json:"principal"`
	TermMonths int     `json:"term_months"`
}

func (r *CreateCreditRequest) Validate() error {
	if r.AccountID == "" {
		return NewErrorResponse("Account ID is required")
	}
	if r.Principal <= 0 {
		return NewErrorResponse("Principal must be greater than zero")
	}
	if r.TermMonths <= 0 {
		return NewErrorResponse("Term must be greater than zero months")
	}

	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if r.Email == "" {
		return NewErrorResponse("Email is required")
	} else if strings.Contains(r.Email, "@") == false {
		return NewErrorResponse("Email must be valid")
	}
	if r.Password == "" {
		return NewErrorResponse("Password is required")
	} else if len(r.Password) < 8 {
		return NewErrorResponse("Password must be at least 8 characters long")
	}

	return nil
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Validate() error {
	if r.Username == "" {
		return NewErrorResponse("Username is required")
	}
	if r.Email == "" {
		return NewErrorResponse("Email is required")
	} else if strings.Contains(r.Email, "@") == false {
		return NewErrorResponse("Email must be valid")
	}
	if r.Password == "" {
		return NewErrorResponse("Password is required")
	} else if len(r.Password) < 8 {
		return NewErrorResponse("Password must be at least 8 characters long")
	}

	return nil
}

type TransferRequest struct {
	FromAccountID string  `json:"from_account_id"`
	ToAccountID   string  `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

func (r *TransferRequest) Validate() error {
	if r.FromAccountID == "" {
		return NewErrorResponse("From Account ID is required")
	}
	if r.ToAccountID == "" {
		return NewErrorResponse("To Account ID is required")
	}
	if r.Amount <= 0 {
		return NewErrorResponse("Amount must be greater than zero")
	}

	return nil
}
