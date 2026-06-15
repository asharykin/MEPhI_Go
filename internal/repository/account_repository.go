package repository

import (
	"banksystem/internal/model"
	"context"
	"database/sql"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, tx *sql.Tx, account *model.Account) (*model.Account, error) {
	query := `
		INSERT INTO accounts (user_id, balance, type, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := tx.QueryRowContext(
		ctx,
		query,
		account.UserID,
		account.Balance,
		account.Type,
		account.IsActive,
		account.CreatedAt,
		account.UpdatedAt,
	).Scan(&account.ID, &account.CreatedAt)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *AccountRepository) GetByID(ctx context.Context, id int64) (*model.Account, error) {
	query := `
		SELECT id, user_id, balance, type, is_active, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`

	account := &model.Account{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.Balance,
		&account.Type,
		&account.IsActive,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return account, err
}

func (r *AccountRepository) GetByUserID(ctx context.Context, userID int64) ([]*model.Account, error) {
	query := `
		SELECT id, user_id, balance, type, is_active, created_at, updated_at
		FROM accounts
		WHERE user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*model.Account
	for rows.Next() {
		account := &model.Account{}
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Balance,
			&account.Type,
			&account.IsActive,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, rows.Err()
}

func (r *AccountRepository) Update(ctx context.Context, tx *sql.Tx, account *model.Account) error {
	query := `
		UPDATE accounts
		SET balance = $1,
			is_active = $2,
			updated_at = $3
		WHERE id = $4
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		account.Balance,
		account.IsActive,
		account.UpdatedAt,
		account.ID,
	)

	return err
}

func (r *AccountRepository) UpdateBalance(tx *sql.Tx, id int64, balance float64) error {
	query := `
		UPDATE accounts
		SET balance = $1
		WHERE id = $2
	`

	_, err := tx.Exec(query, balance, id)
	return err
}
