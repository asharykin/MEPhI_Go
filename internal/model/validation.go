package model

import (
	"regexp"
	"strings"
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,50}$`)
)

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ValidateUsername(username string) bool {
	return usernameRegex.MatchString(username)
}

func ValidatePassword(password string) bool {
	return len(password) >= 8
}

func ValidateAmount(amount float64) bool {
	return amount > 0
}

func ValidateCardNumber(number string) bool {
	number = strings.ReplaceAll(number, " ", "")
	return len(number) == 16
}

func ValidateCVV(cvv string) bool {
	return len(cvv) == 3
}

func ValidateExpiryDate(date string) bool {
	parts := strings.Split(date, "/")
	if len(parts) != 2 {
		return false
	}

	month := parts[0]
	year := parts[1]

	return len(month) == 2 && len(year) == 2
}
