package dto

import "strings"

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
