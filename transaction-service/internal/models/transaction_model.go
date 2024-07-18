package models

import "time"

type Transaction struct {
	ID        string    `json:"id"`
	SKU       string    `json:"sku"`
	Amount    int32     `json:"amount"`
	Qty       int32     `json:"qty"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTransactionResponse struct {
	ID string `json:"id"`
}

type GetAllTransactionRequest struct {
	Limit  int `json:"limit"`
	Page   int `json:"page"`
	Offset int
}
