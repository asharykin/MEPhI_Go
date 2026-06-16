package repository

import (
	"banksystem/internal/model"
	"database/sql"
	"time"
)

type CreditRepository struct {
	db *sql.DB
}

func NewCreditRepository(db *sql.DB) *CreditRepository {
	return &CreditRepository{
		db: db,
	}
}

func (r *CreditRepository) Create(credit *model.Credit) error {
	query := `
		INSERT INTO credits (user_id, amount, term_months, interest_rate, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		credit.UserID,
		credit.Amount,
		credit.TermMonths,
		credit.InterestRate,
		credit.Status,
		time.Now(),
	).Scan(&credit.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CreditRepository) GetByID(id int) (*model.Credit, error) {
	query := `
		SELECT id, user_id, amount, term_months, interest_rate, status, created_at
		FROM credits
		WHERE id = $1
	`
	credit := &model.Credit{}
	err := r.db.QueryRow(query, id).Scan(
		&credit.ID,
		&credit.UserID,
		&credit.Amount,
		&credit.TermMonths,
		&credit.InterestRate,
		&credit.Status,
		&credit.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return credit, nil
}

func (r *CreditRepository) GetByUserID(userID int) ([]*model.Credit, error) {
	query := `
		SELECT id, user_id, amount, term_months, interest_rate, status, created_at
		FROM credits
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credits []*model.Credit
	for rows.Next() {
		credit := &model.Credit{}
		err := rows.Scan(
			&credit.ID,
			&credit.UserID,
			&credit.Amount,
			&credit.TermMonths,
			&credit.InterestRate,
			&credit.Status,
			&credit.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		credits = append(credits, credit)
	}
	return credits, nil
}

func (r *CreditRepository) UpdateStatus(id int, status string) error {
	query := `
		UPDATE credits
		SET status = $1
		WHERE id = $2
	`
	_, err := r.db.Exec(query, status, id)
	if err != nil {
		return err
	}
	return nil
}
