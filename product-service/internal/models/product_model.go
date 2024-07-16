package models

import "time"

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	SKU         string    `json:"sku"`
	Image       string    `json:"image"`
	Price       int64     `json:"price"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"username"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	Name        string `json:"name"`
	SKU         string `json:"sku"`
	Image       string `json:"image"`
	Price       int64  `json:"price"`
	Description string `json:"description"`
	CreatedBy   string `json:"username"`
}

type CreateProductResponse struct {
	ID int64 `json:"id"`
}

type UpdateProductRequest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	SKU         string `json:"sku"`
	Image       string `json:"image"`
	Price       int64  `json:"price"`
	Description string `json:"description"`
}

type GetAllProductRequest struct {
	Limit  int `json:"limit"`
	Page   int `json:"page"`
	Offset int
}
