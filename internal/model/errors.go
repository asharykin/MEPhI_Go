package model

import "errors"

var (
	// Общие ошибки
	ErrInvalidID = errors.New("неверный ID")

	// Ошибки пользователя
	ErrInvalidUserID      = errors.New("неверный ID пользователя")
	ErrInvalidUsername    = errors.New("неверное имя пользователя")
	ErrInvalidEmail       = errors.New("неверный email")
	ErrInvalidPassword    = errors.New("неверный пароль")
	ErrEmailAlreadyExists = errors.New("email уже существует")
	ErrUserNotFound       = errors.New("пользователь не найден")
	ErrInvalidCredentials = errors.New("неверные учетные данные")

	// Ошибки аккаунта
	ErrInvalidAccountID   = errors.New("неверный ID счета")
	ErrInvalidAccountType = errors.New("неверный тип счета")
	ErrInsufficientFunds  = errors.New("недостаточно средств")
	ErrAccountNotFound    = errors.New("счет не найден")

	// Ошибки карты
	ErrInvalidCardNumber = errors.New("неверный номер карты")
	ErrInvalidExpiryDate = errors.New("неверный срок действия")
	ErrInvalidCVV        = errors.New("неверный CVV")
	ErrCardNotFound      = errors.New("карта не найдена")

	// Ошибки кредита
	ErrInvalidCreditID     = errors.New("неверный ID кредита")
	ErrInvalidAmount       = errors.New("неверная сумма")
	ErrInvalidTerm         = errors.New("неверный срок")
	ErrInvalidRate         = errors.New("неверная процентная ставка")
	ErrInvalidInterestRate = errors.New("неверная процентная ставка")
	ErrCreditNotFound      = errors.New("кредит не найден")

	// Ошибки платежа
	ErrInvalidPaymentID   = errors.New("неверный ID платежа")
	ErrInvalidPaymentDate = errors.New("неверная дата платежа")
	ErrPaymentNotFound    = errors.New("платеж не найден")

	// Ошибки репозитория
	ErrNotFound      = errors.New("запись не найдена")
	ErrAlreadyExists = errors.New("запись уже существует")
)
