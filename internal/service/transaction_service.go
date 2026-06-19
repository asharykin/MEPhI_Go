package service

import (
	"banksystem/internal/dto"
	"banksystem/internal/logger"
	"banksystem/internal/model"
	"banksystem/internal/repository"
	"banksystem/internal/util"
	"context"
	"errors"
	"fmt"
	"time"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	accountRepo     *repository.AccountRepository
	storage         *repository.Storage
}

func NewTransactionService(transactionRepo *repository.TransactionRepository, accountRepo *repository.AccountRepository, storage *repository.Storage) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		storage:         storage,
	}
}

func (s *TransactionService) Transfer(ctx context.Context, userID string, req *dto.TransferRequest) error {
	fromAccount, err := s.accountRepo.GetByIDAndUserID(ctx, req.FromAccountID, userID)
	if err != nil {
		logger.Error("Failed to get sender account for transfer", "error", err, "account_id", req.FromAccountID, "user_id", userID)
		return err
	}

	toAccount, err := s.accountRepo.GetByID(ctx, req.ToAccountID)
	if err != nil {
		logger.Error("Failed to get receiver account for transfer", "error", err, "account_id", req.ToAccountID)
		return fmt.Errorf("failed to get receiver account: %w", err)
	}
	if toAccount == nil {
		return errors.New("receiver account not found")
	}

	if fromAccount.Balance < req.Amount {
		logger.Warn("Insufficient funds for transfer", "account_id", req.FromAccountID, "balance", fromAccount.Balance, "amount", req.Amount)
		return errors.New("insufficient funds")
	}

	tx, err := s.storage.DB.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Failed to begin transaction for transfer", "error", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	err = s.accountRepo.UpdateBalanceTx(ctx, tx, req.FromAccountID, fromAccount.Balance-req.Amount)
	if err != nil {
		logger.Error("Failed to update sender balance for transfer", "error", err, "account_id", req.FromAccountID)
		return fmt.Errorf("failed to update sender balance: %w", err)
	}

	err = s.accountRepo.UpdateBalanceTx(ctx, tx, req.ToAccountID, toAccount.Balance+req.Amount)
	if err != nil {
		logger.Error("Failed to update receiver balance for transfer", "error", err, "account_id", req.ToAccountID)
		return fmt.Errorf("failed to update receiver balance: %w", err)
	}

	transaction := &model.Transaction{
		ID:          util.GenerateUUID(),
		SenderID:    req.FromAccountID,
		ReceiverID:  req.ToAccountID,
		Amount:      req.Amount,
		Type:        "transfer",
		Description: "Transfer between accounts",
		CreatedAt:   time.Now(),
	}

	err = s.transactionRepo.CreateTx(ctx, tx, transaction)
	if err != nil {
		logger.Error("Failed to create transaction record for transfer", "error", err, "transaction_id", transaction.ID)
		return fmt.Errorf("failed to create transaction record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		logger.Error("Failed to commit transaction for transfer", "error", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info("Transfer successful", "from_account", req.FromAccountID, "to_account", req.ToAccountID, "amount", req.Amount)
	return nil
}

func (s *TransactionService) GetTransactionsBySender(ctx context.Context, senderID string) ([]*model.Transaction, error) {
	return s.transactionRepo.GetBySenderID(ctx, senderID)
}

func (s *TransactionService) GetTransactionsByReceiver(ctx context.Context, receiverID string) ([]*model.Transaction, error) {
	return s.transactionRepo.GetByReceiverID(ctx, receiverID)
}
