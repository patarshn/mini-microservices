package repository

import (
	"database/sql"
	"log"
	"time"
	"transaction-service/internal/models"
)

func GetTransactionData(db *sql.DB, thresholdDate time.Time) ([]models.Transaction, error) {
	var (
		transactions []models.Transaction
	)
	rows, err := db.Query("SELECT id, sku, amount, qty, created_by, created_at, updated_at FROM transactions WHERE updated_at < ?", thresholdDate)

	if err != nil {
		log.Fatalf("Failed to fetch records: %v", err)
		return transactions, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Transaction
		if err := rows.Scan(&r.ID, &r.SKU, &r.Amount, &r.Qty, &r.CreatedBy, &r.CreatedAt, &r.UpdatedAt); err != nil {
			log.Fatalf("Failed to scan record: %v", err)
			return transactions, err
		}
		transactions = append(transactions, r)
	}

	return transactions, nil
}

func DeleteTransaction(tx *sql.Tx, thresholdDate time.Time) error {

	query := `ALTER TABLE transactions DELETE WHERE created_at < ?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		thresholdDate,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
