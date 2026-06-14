package repositories

import (
	"banksystem/internal/models"
	"context"
	"database/sql"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *sql.Tx, transaction *models.Transaction) (*models.Transaction, error) {
	query := `
		INSERT INTO transactions (account_id, type, amount, status, to_account_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := tx.QueryRowContext(
		ctx,
		query,
		transaction.AccountID,
		transaction.Type,
		transaction.Amount,
		transaction.Status,
		transaction.ToAccountID,
		transaction.CreatedAt,
	).Scan(&transaction.ID, &transaction.CreatedAt)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, id int64) (*models.Transaction, error) {
	query := `
		SELECT id, account_id, type, amount, status, to_account_id, created_at
		FROM transactions
		WHERE id = $1
	`

	transaction := &models.Transaction{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.AccountID,
		&transaction.Type,
		&transaction.Amount,
		&transaction.Status,
		&transaction.ToAccountID,
		&transaction.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return transaction, err
}

func (r *TransactionRepository) GetByAccountID(ctx context.Context, accountID int64) ([]*models.Transaction, error) {
	query := `
		SELECT id, account_id, type, amount, status, to_account_id, created_at
		FROM transactions
		WHERE account_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		transaction := &models.Transaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.AccountID,
			&transaction.Type,
			&transaction.Amount,
			&transaction.Status,
			&transaction.ToAccountID,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, rows.Err()
}

func (r *TransactionRepository) UpdateStatus(ctx context.Context, tx *sql.Tx, id int64, status string) error {
	query := `
		UPDATE transactions
		SET status = $1
		WHERE id = $2
	`

	_, err := tx.ExecContext(ctx, query, status, id)
	return err
}
