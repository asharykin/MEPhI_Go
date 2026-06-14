package repositories

import (
	"banksystem/internal/models"
	"database/sql"
	"time"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(tx *sql.Tx, payment *models.Payment) error {
	query := `
		INSERT INTO payments (credit_id, amount, due_date, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := tx.QueryRow(
		query,
		payment.CreditID,
		payment.Amount,
		payment.DueDate,
		payment.Status,
		time.Now(),
	).Scan(&payment.ID)

	return err
}

func (r *PaymentRepository) GetByID(id int64) (*models.Payment, error) {
	query := `
		SELECT id, credit_id, amount, due_date, status, created_at
		FROM payments
		WHERE id = $1
	`

	payment := &models.Payment{}
	err := r.db.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.CreditID,
		&payment.Amount,
		&payment.DueDate,
		&payment.Status,
		&payment.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return payment, err
}

func (r *PaymentRepository) GetByCreditID(creditID int64) ([]*models.Payment, error) {
	query := `
		SELECT id, credit_id, amount, due_date, status, created_at
		FROM payments
		WHERE credit_id = $1
		ORDER BY due_date ASC
	`

	rows, err := r.db.Query(query, creditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*models.Payment
	for rows.Next() {
		payment := &models.Payment{}
		err := rows.Scan(
			&payment.ID,
			&payment.CreditID,
			&payment.Amount,
			&payment.DueDate,
			&payment.Status,
			&payment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, rows.Err()
}

func (r *PaymentRepository) GetPendingPayments() ([]*models.Payment, error) {
	query := `
		SELECT id, credit_id, amount, due_date, status, created_at
		FROM payments
		WHERE status = 'pending' AND due_date <= $1
		ORDER BY due_date ASC
	`

	rows, err := r.db.Query(query, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*models.Payment
	for rows.Next() {
		payment := &models.Payment{}
		err := rows.Scan(
			&payment.ID,
			&payment.CreditID,
			&payment.Amount,
			&payment.DueDate,
			&payment.Status,
			&payment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, rows.Err()
}

func (r *PaymentRepository) UpdateStatus(tx *sql.Tx, id int64, status string) error {
	query := `
		UPDATE payments
		SET status = $1
		WHERE id = $2
	`

	_, err := tx.Exec(query, status, id)
	return err
}
