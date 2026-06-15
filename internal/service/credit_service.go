package service

import (
	"banksystem/internal/model"
	"banksystem/internal/repository"
	"context"
	"database/sql"
	"errors"
	"math"
	"time"
)

type CreditService struct {
	creditRepo      *repository.CreditRepository
	accountRepo     *repository.AccountRepository
	paymentRepo     *repository.CreditPaymentRepository
	transactionRepo *repository.TransactionRepository
	db              *sql.DB
}

func NewCreditService(
	db *sql.DB,
	creditRepo *repository.CreditRepository,
	accountRepo *repository.AccountRepository,
	paymentRepo *repository.CreditPaymentRepository,
	transactionRepo *repository.TransactionRepository,
) *CreditService {
	return &CreditService{
		db:              db,
		creditRepo:      creditRepo,
		accountRepo:     accountRepo,
		paymentRepo:     paymentRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *CreditService) CreateCredit(ctx context.Context, userID int64, accountID int64, amount float64, term int, rate float64) (*model.Credit, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("счет не найден")
	}

	monthlyPayment := s.calculateMonthlyPayment(amount, rate, term)

	credit := &model.Credit{
		UserID:         userID,
		AccountID:      accountID,
		Amount:         amount,
		InterestRate:   rate,
		TermMonths:     term,
		MonthlyPayment: monthlyPayment,
		Status:         model.CreditStatusActive,
		CreatedAt:      time.Now(),
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = s.creditRepo.Create(credit)
	if err != nil {
		return nil, err
	}

	err = s.createPaymentSchedule(ctx, credit)
	if err != nil {
		return nil, err
	}

	account.Balance += amount
	err = s.accountRepo.Update(ctx, tx, account)
	if err != nil {
		return nil, err
	}

	transaction := &model.Transaction{
		AccountID: accountID,
		Type:      "credit",
		Amount:    amount,
		Status:    "completed",
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	_, err = s.transactionRepo.Create(ctx, tx, transaction)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return credit, nil
}

func (s *CreditService) GetCredit(ctx context.Context, creditID int64) (*model.Credit, error) {
	return s.creditRepo.GetByID(int(creditID))
}

func (s *CreditService) GetUserCredits(ctx context.Context, userID int64) ([]*model.Credit, error) {
	return s.creditRepo.GetByUserID(int(userID))
}

func (s *CreditService) GetPaymentSchedule(ctx context.Context, creditID int64) ([]*model.CreditPayment, error) {
	return s.paymentRepo.GetByCreditID(ctx, creditID)
}

func (s *CreditService) calculateMonthlyPayment(amount float64, rate float64, term int) float64 {
	monthlyRate := rate / 12 / 100
	return amount * monthlyRate * math.Pow(1+monthlyRate, float64(term)) / (math.Pow(1+monthlyRate, float64(term)) - 1)
}

func (s *CreditService) createPaymentSchedule(ctx context.Context, credit *model.Credit) error {
	for i := 1; i <= credit.TermMonths; i++ {
		payment := &model.CreditPayment{
			CreditID: credit.ID,
			Amount:   credit.MonthlyPayment,
			Status:   "pending",
			DueDate:  credit.CreatedAt.AddDate(0, i, 0),
		}

		err := s.paymentRepo.Create(ctx, payment)
		if err != nil {
			return err
		}
	}

	return nil
}
