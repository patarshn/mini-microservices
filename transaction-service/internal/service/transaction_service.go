package service

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"transaction-service/internal/external/product"
	"transaction-service/internal/models"
	"transaction-service/internal/repository"

	"github.com/google/uuid"
)

type ITransactionService interface {
	CreateTransaction(ctx context.Context, transaction models.Transaction) (string, error)
	UpdateTransaction(ctx context.Context, transaction models.Transaction) error
	GetTransactionByID(ctx context.Context, id string) (models.Transaction, error)
	DeleteTransaction(ctx context.Context, id string) error
	GetAllTransaction(ctx context.Context, req models.GetAllTransactionRequest) (models.Paginate, error)
}

type TransactionService struct {
	repo        repository.ITransactionRepository
	productRepo product.IProductRepository
}

func NewTransactionService(
	repo repository.ITransactionRepository,
	productRepo product.IProductRepository,
) ITransactionService {
	return &TransactionService{
		repo:        repo,
		productRepo: productRepo}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, transaction models.Transaction) (string, error) {

	var (
		product *product.Product
	)

	product, err := s.productRepo.GetProductBySKU(ctx, transaction.SKU)
	if err != nil {
		return "", err
	}

	if product == nil {
		return "", errors.New("SKU Not valid")
	}

	transaction.ID = uuid.New().String()
	existTransactionID, err := s.repo.GetTransactionByID(ctx, transaction.ID)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return "", err
	}

	if existTransactionID.ID == transaction.ID {
		return "", errors.New("ID Already Used")
	}

	transaction.Amount = int32(product.Price) * transaction.Qty
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()
	transaction.CreatedBy = "TESTING"

	err = s.repo.CreateTransaction(ctx, transaction)
	if err != nil {
		return "", err
	}

	return transaction.ID, nil
}

func (s *TransactionService) UpdateTransaction(ctx context.Context, transaction models.Transaction) error {

	var (
		product *product.Product
	)

	product, err := s.productRepo.GetProductBySKU(ctx, transaction.SKU)
	if err != nil {
		return err
	}

	if product == nil {
		return errors.New("SKU Not valid")
	}

	transaction.Amount = int32(product.Price) * transaction.Qty
	transaction.UpdatedAt = time.Now()

	return s.repo.UpdateTransaction(ctx, transaction)
}

func (s *TransactionService) GetTransactionByID(ctx context.Context, id string) (models.Transaction, error) {
	return s.repo.GetTransactionByID(ctx, id)
}

func (s *TransactionService) DeleteTransaction(ctx context.Context, id string) error {
	return s.repo.DeleteTransaction(ctx, id)
}

func (s *TransactionService) GetAllTransaction(ctx context.Context, req models.GetAllTransactionRequest) (models.Paginate, error) {
	var paginateRes models.Paginate

	req.Offset = (req.Page * req.Limit) - req.Limit

	transactions, err := s.repo.GetAllTransaction(ctx, req)
	if err != nil {
		return models.Paginate{}, err
	}

	paginateRes.Next = false
	loopLimit := len(transactions)
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
		paginateRes.Data = append(paginateRes.Data, transactions[i])
	}

	paginateRes.From = req.Offset + 1
	paginateRes.To = req.Offset + req.Limit
	paginateRes.Page = int64(req.Page)

	return paginateRes, nil
}
