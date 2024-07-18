package usecase

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"transaction-service/cmd/cron/repository"

	"github.com/minio/minio-go/v7"
)

func RunMigrationJob(db *sql.DB, minioClient *minio.Client) {
	fmt.Println("RunMigrationJob -- Start")

	thresholdDate := time.Now().AddDate(0, 0, 0)

	records, err := repository.GetTransactionData(db, thresholdDate)
	if err != nil {
		log.Fatalf("Failed to getTransactionData: %v", err)
	}
	var csvBuffer bytes.Buffer
	for _, r := range records {
		csvBuffer.WriteString(fmt.Sprintf("%s,%s,%d,%d,%s,%s,%s\n", r.ID, r.SKU, r.Amount, r.Qty, r.CreatedBy, r.CreatedAt.Format(time.RFC3339), r.UpdatedAt.Format(time.RFC3339)))
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to initialize tx: %v", err)
	}

	objectName := fmt.Sprintf("transactions_cold_%s.csv", time.Now().Format("2006-01-02"))
	reader := bytes.NewReader(csvBuffer.Bytes())
	size := int64(reader.Len())
	_, err = minioClient.PutObject(context.Background(), "transaction", objectName, reader, size, minio.PutObjectOptions{ContentType: "text/csv"})
	if err != nil {
		log.Fatalf("Failed to upload to MinIO: %v", err)
	}
	log.Printf("Successfully uploaded %s to MinIO", objectName)

	err = repository.DeleteTransaction(tx, thresholdDate)
	if err != nil {
		log.Fatalf("Failed to deleteTransaction: %v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit tx: %v", err)
	}

	log.Println("Successfully deleted old records from ClickHouse")
	fmt.Println("RunMigrationJob -- End")
}
