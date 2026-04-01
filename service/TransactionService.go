package service

import (
	"fmt"
	"go-back/repository"
)

type TransactionService struct {
	txRepo      *repository.TransactionRepository
	cardRepo    *repository.CardRepository
	authService *AuthService
}

type AuthorizationResponse struct {
	Authorized bool   `json:"authorized"`
	Message    string `json:"message"`
	CardNumber string `json:"card_number,omitempty"`
	Balance    float64 `json:"balance,omitempty"`
}

func NewTransactionService(txRepo *repository.TransactionRepository, cardRepo *repository.CardRepository, authService *AuthService) *TransactionService {
	return &TransactionService{
		txRepo:      txRepo,
		cardRepo:    cardRepo,
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

func (s *TransactionService) AuthorizeTransaction(cardNumber string, amount float64) *AuthorizationResponse {
	// Проверка входных данных
	if cardNumber == "" {
		return &AuthorizationResponse{
			Authorized: false,
			Message:    "Card number cannot be empty",
		}
	}

	if amount <= 0 {
		return &AuthorizationResponse{
			Authorized: false,
			Message:    "Transaction amount must be greater than 0",
		}
	}

	cards, err := s.cardRepo.GetAllCards()
	if err != nil {
		return &AuthorizationResponse{
			Authorized: false,
			Message:    "Failed to retrieve card information",
		}
	}

	var targetCard *repository.Card
	for i := range cards {
		if cards[i].CardNumber == cardNumber {
			targetCard = &cards[i]
			break
		}
	}

	if targetCard == nil {
		return &AuthorizationResponse{
			Authorized: false,
			Message:    "Card not found",
		}
	}

	if targetCard.IsBlocked == 1 {
		return &AuthorizationResponse{
			Authorized: false,
			Message:    "Card is blocked",
			CardNumber: cardNumber,
		}
	}

	if targetCard.Balance < amount {
		return &AuthorizationResponse{
			Authorized: false,
			Message:    fmt.Sprintf("Insufficient funds. Required: %.2f, Available: %.2f", amount, targetCard.Balance),
			CardNumber: cardNumber,
			Balance:    targetCard.Balance,
		}
	}

	return &AuthorizationResponse{
		Authorized: true,
		Message:    "Transaction authorized successfully",
		CardNumber: cardNumber,
		Balance:    targetCard.Balance,
	}
}