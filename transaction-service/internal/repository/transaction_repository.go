package repository

import (
	"context"
	"database/sql"
	"fmt"
	"transaction-service/internal/models"
)

type ITransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction models.Transaction) error
	UpdateTransaction(ctx context.Context, transaction models.Transaction) error
	GetTransactionByID(ctx context.Context, id string) (models.Transaction, error)
	DeleteTransaction(ctx context.Context, id string) error
	GetAllTransaction(ctx context.Context, req models.GetAllTransactionRequest) ([]models.Transaction, error)
}

type TransactionRepository struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) ITransactionRepository {
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction models.Transaction) error {

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	query := `INSERT INTO transactions (id, sku, amount, qty, created_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		transaction.ID,
		transaction.SKU,
		transaction.Amount,
		transaction.Qty,
		transaction.CreatedBy,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *TransactionRepository) UpdateTransaction(ctx context.Context, transaction models.Transaction) error {

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	query := `
	ALTER TABLE transactions
		UPDATE sku=?, amount=?, qty=?, updated_at=? 
	WHERE id=?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		transaction.SKU,
		transaction.Amount,
		transaction.Qty,
		transaction.UpdatedAt,
		transaction.ID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *TransactionRepository) GetTransactionByID(ctx context.Context, id string) (models.Transaction, error) {
	query := `SELECT id, sku, amount, qty, created_by, created_at, updated_at FROM transactions WHERE id=?`
	transaction := models.Transaction{}
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.SKU,
		&transaction.Amount,
		&transaction.Qty,
		&transaction.CreatedBy,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (r *TransactionRepository) DeleteTransaction(ctx context.Context, id string) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	query := `
	ALTER TABLE transactions DELETE WHERE id=?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		id,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *TransactionRepository) GetAllTransaction(ctx context.Context, req models.GetAllTransactionRequest) ([]models.Transaction, error) {
	query := fmt.Sprintf(`
		SELECT id, sku, amount, qty, created_by, created_at, updated_at
		FROM transactions
		ORDER BY updated_at DESC
		LIMIT %d OFFSET %d`,
		req.Limit,
		req.Offset,
	)
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		transaction := models.Transaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.SKU,
			&transaction.Amount,
			&transaction.Qty,
			&transaction.CreatedBy,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
