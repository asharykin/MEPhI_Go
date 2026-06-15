package repository

import (
	"database/sql"
	"errors"
)

var (
	ErrNotFound      = errors.New("запись не найдена")
	ErrAlreadyExists = errors.New("запись уже существует")
	ErrInvalidData   = errors.New("недопустимые данные")
)

type BaseRepository struct {
	db *sql.DB
}

func NewBaseRepository(db *sql.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

func (r *BaseRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *BaseRepository) RollbackTx(tx *sql.Tx) {
	if tx != nil {
		tx.Rollback()
	}
}

func (r *BaseRepository) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}
