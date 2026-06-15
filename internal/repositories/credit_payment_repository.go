package repositories

import (
	"banksystem/internal/model"
	"context"
	"database/sql"
	"time"
)

type CreditPaymentRepository struct {
	db *sql.DB
}

func NewCreditPaymentRepository(db *sql.DB) *CreditPaymentRepository {
	return &CreditPaymentRepository{db: db}
}

func (r *CreditPaymentRepository) Create(ctx context.Context, payment *model.CreditPayment) error {
	query := `
		INSERT INTO credit_payments (credit_id, amount, status, due_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(ctx, query,
		payment.CreditID,
		payment.Amount,
		payment.Status,
		payment.DueDate,
	).Scan(&payment.ID, &payment.CreatedAt)
}

func (r *CreditPaymentRepository) GetByID(ctx context.Context, id int64) (*model.CreditPayment, error) {
	query := `
		SELECT id, credit_id, amount, status, due_date, created_at
		FROM credit_payments
		WHERE id = $1
	`

	payment := &model.CreditPayment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.CreditID,
		&payment.Amount,
		&payment.Status,
		&payment.DueDate,
		&payment.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return payment, err
}

func (r *CreditPaymentRepository) GetByCreditID(ctx context.Context, creditID int64) ([]*model.CreditPayment, error) {
	query := `
		SELECT id, credit_id, amount, status, due_date, created_at
		FROM credit_payments
		WHERE credit_id = $1
		ORDER BY due_date
	`

	rows, err := r.db.QueryContext(ctx, query, creditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*model.CreditPayment
	for rows.Next() {
		payment := &model.CreditPayment{}
		err := rows.Scan(
			&payment.ID,
			&payment.CreditID,
			&payment.Amount,
			&payment.Status,
			&payment.DueDate,
			&payment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *CreditPaymentRepository) GetPending(ctx context.Context) ([]*model.CreditPayment, error) {
	query := `
		SELECT id, credit_id, amount, status, due_date, created_at
		FROM credit_payments
		WHERE status = 'pending' AND due_date <= $1
		ORDER BY due_date
	`

	rows, err := r.db.QueryContext(ctx, query, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*model.CreditPayment
	for rows.Next() {
		payment := &model.CreditPayment{}
		err := rows.Scan(
			&payment.ID,
			&payment.CreditID,
			&payment.Amount,
			&payment.Status,
			&payment.DueDate,
			&payment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *CreditPaymentRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `
		UPDATE credit_payments
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}
