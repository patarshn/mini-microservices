package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"transaction-service/internal/external/product"
	"transaction-service/internal/handler"
	"transaction-service/internal/middleware"
	"transaction-service/internal/repository"
	"transaction-service/internal/service"

	"github.com/gorilla/mux"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("clickhouse", os.Getenv("CLICKHOUSE_URI_TCP"))
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	fmt.Println("Init Repository...")
	repo := repository.NewTransactionRepository(db)
	productRepo := product.NewProductRepository(os.Getenv("PRODUCT_SERVICE_URI"), httpClient)

	fmt.Println("Init Repository...")
	transactionService := service.NewTransactionService(repo, productRepo)

	fmt.Println("Init Handler...")
	transactionHandler := handler.NewTransactionHandler(transactionService)

	r := mux.NewRouter()

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}).Methods("GET")

	r.Use(middleware.AuthMiddleware)

	r.HandleFunc("/transactions", transactionHandler.GetAllTransaction).Methods(http.MethodGet)
	r.HandleFunc("/transactions/{id}", transactionHandler.GetTransactionByID).Methods(http.MethodGet)
	r.HandleFunc("/transactions", transactionHandler.CreateTransaction).Methods(http.MethodPost)
	r.HandleFunc("/transactions/{id}", transactionHandler.UpdateTransaction).Methods(http.MethodPut)
	r.HandleFunc("/transactions/{id}", transactionHandler.DeleteTransaction).Methods(http.MethodDelete)

	//TODO: port better use .env
	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8083",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Starting server on :8083")
	log.Fatal(srv.ListenAndServe())
}
