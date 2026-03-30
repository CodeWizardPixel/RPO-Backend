package repository

import (
	"database/sql"
	"fmt"
	"time"
)

type TransactionRepository struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

type Transaction struct {
	ID         int
	Amount     float64
	CardID     int
	TerminalID int
	CreatedAt  time.Time
}

func (r *TransactionRepository) GetAllTransactions() ([]Transaction, error) {
	rows, err := r.DB.Query("select * from transactions")
	if err != nil {
		return nil, fmt.Errorf("error querying transactions: %w", err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.Amount, &t.CardID, &t.TerminalID, &t.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning transaction row: %w", err)
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetTransactionByID(id int) (*Transaction, error) {
	row := r.DB.QueryRow("select * from transactions where id = ?", id)
	var t Transaction
	err := row.Scan(&t.ID, &t.Amount, &t.CardID, &t.TerminalID, &t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error scanning transaction row: %w", err)
	}
	return &t, nil
}

func (r *TransactionRepository) CreateTransaction(amount float64, cardID, terminalID int) error {
	_, err := r.DB.Exec("insert into transactions (amount, card_id, terminal_id) values (?, ?, ?)",
		amount, cardID, terminalID)
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepository) DeleteTransactionByID(id int) error {
	_, err := r.DB.Exec("delete from transactions where id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting transaction: %w", err)
	}
	return nil
}
