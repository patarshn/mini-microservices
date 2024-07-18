package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"transaction-service/cmd/cron/usecase"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	cron "github.com/robfig/cron/v3"
)

func main() {
	log.Println("Cron main() Start")

	log.Println("Init env...")
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	log.Println("Init click...")
	db, err := sql.Open("clickhouse", os.Getenv("CLICKHOUSE_URI_TCP"))
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()

	log.Println("Init minio...")
	minioClient, err := minio.New(os.Getenv("MINIO_HOST"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESS"), os.Getenv("MINIO_SECRET"), ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	log.Println("Init schedular...")
	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	scheduler := cron.New(cron.WithLocation(jakartaTime))
	defer scheduler.Stop()

	log.Println("Run schedular...")
	scheduler.AddFunc("0 0 */1 * *", func() { usecase.RunMigrationJob(db, minioClient) })

	go scheduler.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

}
