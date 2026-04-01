package service

import (
	"fmt"
	"go-back/repository"
)

type TransactionService struct {
	txRepo      *repository.TransactionRepository
	authService *AuthService
}

func NewTransactionService(txRepo *repository.TransactionRepository, authService *AuthService) *TransactionService {
	return &TransactionService{
		txRepo:      txRepo,
		authService: authService,
	}
}

func (s *TransactionService) GetAllTransactions() ([]repository.Transaction, error) {
	txs, err := s.txRepo.GetAllTransactions()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all transactions: %w", err)
	}
	return txs, nil
}

func (s *TransactionService) GetTransactionByID(id int) (*repository.Transaction, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid transaction ID: must be greater than 0")
	}

	tx, err := s.txRepo.GetTransactionByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transaction with ID %d: %w", id, err)
	}
	return tx, nil
}

func (s *TransactionService) CreateTransaction(tokenString string, amount float64, cardID, terminalID int) error {
	if err := ensureAdmin(s.authService, tokenString); err != nil {
		return err
	}
	if cardID <= 0 {
		return fmt.Errorf("invalid card ID: must be greater than 0")
	}
	if terminalID <= 0 {
		return fmt.Errorf("invalid terminal ID: must be greater than 0")
	}

	if err := s.txRepo.CreateTransaction(amount, cardID, terminalID); err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

func (s *TransactionService) DeleteTransaction(tokenString string, id int) error {
	if err := ensureAdmin(s.authService, tokenString); err != nil {
		return err
	}
	if id <= 0 {
		return fmt.Errorf("invalid transaction ID: must be greater than 0")
	}

	if err := s.txRepo.DeleteTransactionByID(id); err != nil {
		return fmt.Errorf("failed to delete transaction with ID %d: %w", id, err)
	}
	return nil
}