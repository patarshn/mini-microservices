package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"product-service/internal/handler"
	"product-service/internal/middleware"
	"product-service/internal/repository"
	"product-service/internal/service"
	"time"

	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URI"))
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	defer db.Close()

	fmt.Println("Init Repository...")
	repo := repository.NewProductRepository(db)

	fmt.Println("Init Repository...")
	productService := service.NewProductService(repo)

	fmt.Println("Init Handler...")
	productHandler := handler.NewProductHandler(productService)

	r := mux.NewRouter()

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}).Methods("GET")

	r.Use(middleware.AuthMiddleware)

	r.HandleFunc("/products", productHandler.GetAllProduct).Methods(http.MethodGet)
	r.HandleFunc("/products/{id:[0-9]+}", productHandler.GetProductByID).Methods(http.MethodGet)
	r.HandleFunc("/products/sku/{sku}", productHandler.GetProductBySKU).Methods(http.MethodGet)
	r.HandleFunc("/products", productHandler.CreateProduct).Methods(http.MethodPost)
	r.HandleFunc("/products/{id:[0-9]+}", productHandler.UpdateProduct).Methods(http.MethodPut)
	r.HandleFunc("/products/{id:[0-9]+}", productHandler.DeleteProduct).Methods(http.MethodDelete)

	//TODO: port better use .env
	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8082",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Starting server on :8082")
	log.Fatal(srv.ListenAndServe())
}
