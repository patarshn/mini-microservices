package repository

import (
	"context"
	"database/sql"
	"errors"
	"product-service/internal/models"
)

type IProductRepository interface {
	CreateProduct(ctx context.Context, product models.Product) (int64, error)
	UpdateProduct(ctx context.Context, product models.Product) error
	GetProductByID(ctx context.Context, id int64) (models.Product, error)
	GetProductBySKU(ctx context.Context, sku string) (models.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
	GetAllProduct(ctx context.Context, req models.GetAllProductRequest) ([]models.Product, error)
}

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) IProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product models.Product) (int64, error) {
	var id int64

	query := `
		INSERT INTO products (name, sku, image, price, description, created_by, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	err := r.DB.QueryRowContext(ctx, query,
		product.Name,
		product.SKU,
		product.Image,
		product.Price,
		product.Description,
		product.CreatedBy,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&id)

	return id, err
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product models.Product) error {
	var existID *int

	query := `UPDATE products SET name=$1, sku=$2, image=$3, price=$4, description=$5, updated_at=$6 WHERE id=$7 RETURNING id`
	err := r.DB.QueryRowContext(ctx, query,
		product.Name,
		product.SKU,
		product.Image,
		product.Price,
		product.Description,
		product.UpdatedAt,
		product.ID,
	).Scan(&existID)

	if existID == nil {
		return errors.New("Product ID Not found")
	}
	return err
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id int64) (models.Product, error) {
	query := `SELECT id, name, sku, image, price, description, created_by, created_at, updated_at FROM products WHERE id=$1`
	product := models.Product{}
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.SKU,
		&product.Image,
		&product.Price,
		&product.Description,
		&product.CreatedBy,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return product, err
	}
	return product, nil
}

func (r *ProductRepository) GetProductBySKU(ctx context.Context, sku string) (models.Product, error) {
	query := `SELECT id, name, sku, image, price, description, created_by, created_at, updated_at FROM products WHERE sku=$1`
	product := models.Product{}
	err := r.DB.QueryRowContext(ctx, query, sku).Scan(
		&product.ID,
		&product.Name,
		&product.SKU,
		&product.Image,
		&product.Price,
		&product.Description,
		&product.CreatedBy,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return product, err
	}
	return product, nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id=$1 RETURNING id`
	var existID *int
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&existID)
	if existID == nil {
		return errors.New("Product ID Not found")
	}
	return err
}

func (r *ProductRepository) GetAllProduct(ctx context.Context, req models.GetAllProductRequest) ([]models.Product, error) {
	query := `SELECT * FROM products ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.DB.QueryContext(ctx, query, req.Limit+1, req.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		product := models.Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.SKU,
			&product.Image,
			&product.Price,
			&product.Description,
			&product.CreatedBy,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
