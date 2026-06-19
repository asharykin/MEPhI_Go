package repository

import (
	"banksystem/internal/logger"
	"banksystem/internal/model"
	"context"
	"database/sql"
)

type TransactionRepository struct {
	Storage *Storage
}

func NewTransactionRepository(storage *Storage) *TransactionRepository {
	return &TransactionRepository{Storage: storage}
}

func (r *TransactionRepository) Create(ctx context.Context, transaction *model.Transaction) error {
	query := `INSERT INTO transactions (id, sender_id, receiver_id, amount, type, description, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.Storage.DB.ExecContext(ctx, query, transaction.ID, transaction.SenderID, transaction.ReceiverID, transaction.Amount, transaction.Type, transaction.Description, transaction.CreatedAt)
	if err != nil {
		logger.Error("Failed to create transaction in DB", "error", err, "transaction_id", transaction.ID)
		return err
	}
	return nil
}

func (r *TransactionRepository) GetBySenderID(ctx context.Context, senderID string) ([]*model.Transaction, error) {
	rows, err := r.Storage.DB.QueryContext(ctx, `SELECT id, sender_id, receiver_id, amount, type, description, created_at FROM transactions WHERE sender_id = $1`, senderID)
	if err != nil {
		logger.Error("Failed to get transactions by sender ID from DB", "error", err, "sender_id", senderID)
		return nil, err
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.SenderID, &transaction.ReceiverID, &transaction.Amount, &transaction.Type, &transaction.Description, &transaction.CreatedAt); err != nil {
			logger.Error("Failed to scan transaction row", "error", err)
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	if err := rows.Err(); err != nil {
		logger.Error("Error iterating over transaction rows", "error", err)
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) GetByReceiverID(ctx context.Context, receiverID string) ([]*model.Transaction, error) {
	rows, err := r.Storage.DB.QueryContext(ctx, `SELECT id, sender_id, receiver_id, amount, type, description, created_at FROM transactions WHERE receiver_id = $1`, receiverID)
	if err != nil {
		logger.Error("Failed to get transactions by receiver ID from DB", "error", err, "receiver_id", receiverID)
		return nil, err
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.SenderID, &transaction.ReceiverID, &transaction.Amount, &transaction.Type, &transaction.Description, &transaction.CreatedAt); err != nil {
			logger.Error("Failed to scan transaction row", "error", err)
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	if err := rows.Err(); err != nil {
		logger.Error("Error iterating over transaction rows", "error", err)
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) CreateTx(ctx context.Context, tx *sql.Tx, transaction *model.Transaction) error {
	query := `INSERT INTO transactions (id, sender_id, receiver_id, amount, type, description, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := tx.ExecContext(ctx, query, transaction.ID, transaction.SenderID, transaction.ReceiverID, transaction.Amount, transaction.Type, transaction.Description, transaction.CreatedAt)
	if err != nil {
		logger.Error("Failed to create transaction in DB (tx)", "error", err, "transaction_id", transaction.ID)
		return err
	}
	return nil
}
