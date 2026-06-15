package model

import (
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username" validate:"required"`
	Email        string    `json:"email" validate:"required,email"`
	Password     string    `json:"password"` // Только для регистрации
	PasswordHash string    `json:"-"`        // Хеш пароля
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserCreateRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *UserCreateRequest) Validate() error {
	if !ValidateUsername(u.Username) {
		return ErrInvalidUsername
	}
	if !ValidateEmail(u.Email) {
		return ErrInvalidEmail
	}
	if !ValidatePassword(u.Password) {
		return ErrInvalidPassword
	}
	return nil
}

func (u *UserLoginRequest) Validate() error {
	if !ValidateEmail(u.Email) {
		return ErrInvalidEmail
	}
	if u.Password == "" {
		return ErrInvalidPassword
	}
	return nil
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
