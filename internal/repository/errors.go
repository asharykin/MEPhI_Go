package repository

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("запись не найдена")
	ErrAlreadyExists = errors.New("запись уже существует")
)
