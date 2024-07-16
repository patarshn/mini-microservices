package service

import (
	"context"
	"database/sql"
	"errors"
	"product-service/internal/models"
	"product-service/internal/repository"
	"time"
)

type IProductService interface {
	CreateProduct(ctx context.Context, product models.Product) (int64, error)
	UpdateProduct(ctx context.Context, product models.Product) error
	GetProductByID(ctx context.Context, id int64) (models.Product, error)
	GetProductBySKU(ctx context.Context, sku string) (models.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
	GetAllProduct(ctx context.Context, req models.GetAllProductRequest) (models.Paginate, error)
}

type ProductService struct {
	repo repository.IProductRepository
}

func NewProductService(repo repository.IProductRepository) IProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, product models.Product) (int64, error) {

	existProductSKU, err := s.repo.GetProductBySKU(ctx, product.SKU)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return 0, err
	}

	if existProductSKU.SKU == product.SKU {
		return 0, errors.New("SKU Already Used")
	}

	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	product.CreatedBy = "TESTING"

	return s.repo.CreateProduct(ctx, product)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product models.Product) error {
	product.UpdatedAt = time.Now()
	return s.repo.UpdateProduct(ctx, product)
}

func (s *ProductService) GetProductByID(ctx context.Context, id int64) (models.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *ProductService) GetProductBySKU(ctx context.Context, sku string) (models.Product, error) {
	return s.repo.GetProductBySKU(ctx, sku)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.DeleteProduct(ctx, id)
}

func (s *ProductService) GetAllProduct(ctx context.Context, req models.GetAllProductRequest) (models.Paginate, error) {
	var paginateRes models.Paginate

	req.Offset = (req.Page * req.Limit) - req.Limit

	products, err := s.repo.GetAllProduct(ctx, req)
	if err != nil {
		return models.Paginate{}, err
	}

	paginateRes.Next = false
	loopLimit := len(products)
	if loopLimit > req.Limit {
		loopLimit = req.Limit
		paginateRes.Next = true
	}

	paginateRes.Prev = true
	if req.Offset == 0 {
		paginateRes.Prev = false
	}

	paginateRes.Data = []any{}
	for i := 0; i < loopLimit; i++ {
		paginateRes.Data = append(paginateRes.Data, products[i])
	}

	paginateRes.From = req.Offset + 1
	paginateRes.To = req.Offset + req.Limit
	paginateRes.Page = int64(req.Page)

	return paginateRes, nil
}
