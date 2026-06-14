package services

import (
	"banksystem/internal/models"
	"banksystem/internal/repositories"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type AccountService struct {
	accountRepo     *repositories.AccountRepository
	transactionRepo *repositories.TransactionRepository
	userRepo        *repositories.UserRepository
	db              *sql.DB
	smtpService     *SMTPService
}

func NewAccountService(
	db *sql.DB,
	accountRepo *repositories.AccountRepository,
	transactionRepo *repositories.TransactionRepository,
	userRepo *repositories.UserRepository,
	smtpService *SMTPService,
) *AccountService {
	return &AccountService{
		db:              db,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		smtpService:     smtpService,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, userID int64, accountType string) (*models.Account, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	account := &models.Account{
		UserID:    userID,
		Type:      accountType,
		Balance:   0,
		IsActive:  true,
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	account, err = s.accountRepo.Create(ctx, tx, account)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return account, nil
}

func (s *AccountService) GetUserAccounts(ctx context.Context, userID int64) ([]*models.Account, error) {
	accounts, err := s.accountRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user accounts: %v", err)
	}
	return accounts, nil
}

func (s *AccountService) GetBalance(ctx context.Context, accountID int64) (float64, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return 0, fmt.Errorf("failed to get account: %v", err)
	}
	return account.Balance, nil
}

func (s *AccountService) Deposit(ctx context.Context, accountID int64, amount float64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %v", err)
	}

	account.Balance += amount
	account.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	if err := s.accountRepo.Update(ctx, tx, account); err != nil {
		return fmt.Errorf("failed to update account: %v", err)
	}

	transaction := &models.Transaction{
		AccountID: accountID,
		Type:      "deposit",
		Amount:    amount,
		Status:    "completed",
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	if _, err := s.transactionRepo.Create(ctx, tx, transaction); err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	user, err := s.userRepo.GetByID(ctx, account.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	if err := s.smtpService.SendTransactionNotification(user.Email, amount, "deposit"); err != nil {
		return fmt.Errorf("failed to send notification: %v", err)
	}

	return nil
}

func (s *AccountService) Withdraw(ctx context.Context, accountID int64, amount float64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %v", err)
	}

	if account.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}

	account.Balance -= amount
	account.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	if err := s.accountRepo.Update(ctx, tx, account); err != nil {
		return fmt.Errorf("failed to update account: %v", err)
	}

	transaction := &models.Transaction{
		AccountID: accountID,
		Type:      "withdraw",
		Amount:    amount,
		Status:    "completed",
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	if _, err := s.transactionRepo.Create(ctx, tx, transaction); err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	user, err := s.userRepo.GetByID(ctx, account.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	if err := s.smtpService.SendTransactionNotification(user.Email, amount, "withdraw"); err != nil {
		return fmt.Errorf("failed to send notification: %v", err)
	}

	if account.Balance < 1000 {
		if err := s.smtpService.SendLowBalanceNotification(user.Email, account.Balance); err != nil {
			return fmt.Errorf("failed to send low balance notification: %v", err)
		}
	}

	return nil
}

func (s *AccountService) Transfer(ctx context.Context, fromAccountID, toAccountID int64, amount float64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	fromAccount, err := s.accountRepo.GetByID(ctx, fromAccountID)
	if err != nil {
		return fmt.Errorf("failed to get source account: %v", err)
	}

	if fromAccount.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}

	toAccount, err := s.accountRepo.GetByID(ctx, toAccountID)
	if err != nil {
		return fmt.Errorf("failed to get destination account: %v", err)
	}

	fromAccount.Balance -= amount
	toAccount.Balance += amount
	now := time.Now()
	fromAccount.UpdatedAt = sql.NullTime{Time: now, Valid: true}
	toAccount.UpdatedAt = sql.NullTime{Time: now, Valid: true}

	if err := s.accountRepo.Update(ctx, tx, fromAccount); err != nil {
		return fmt.Errorf("failed to update source account: %v", err)
	}

	if err := s.accountRepo.Update(ctx, tx, toAccount); err != nil {
		return fmt.Errorf("failed to update destination account: %v", err)
	}

	transaction := &models.Transaction{
		AccountID:   fromAccountID,
		Type:        "transfer",
		Amount:      amount,
		Status:      "completed",
		ToAccountID: sql.NullInt64{Int64: toAccountID, Valid: true},
		CreatedAt:   sql.NullTime{Time: now, Valid: true},
	}

	if _, err := s.transactionRepo.Create(ctx, tx, transaction); err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	fromUser, err := s.userRepo.GetByID(ctx, fromAccount.UserID)
	if err != nil {
		return fmt.Errorf("failed to get source user: %v", err)
	}

	toUser, err := s.userRepo.GetByID(ctx, toAccount.UserID)
	if err != nil {
		return fmt.Errorf("failed to get destination user: %v", err)
	}

	if err := s.smtpService.SendTransactionNotification(fromUser.Email, amount, "transfer sent"); err != nil {
		return fmt.Errorf("failed to send notification to source user: %v", err)
	}

	if err := s.smtpService.SendTransactionNotification(toUser.Email, amount, "transfer received"); err != nil {
		return fmt.Errorf("failed to send notification to destination user: %v", err)
	}

	if fromAccount.Balance < 1000 {
		if err := s.smtpService.SendLowBalanceNotification(fromUser.Email, fromAccount.Balance); err != nil {
			return fmt.Errorf("failed to send low balance notification: %v", err)
		}
	}

	return nil
}
