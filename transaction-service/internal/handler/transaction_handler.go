package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"transaction-service/internal/models"
	"transaction-service/internal/service"

	"github.com/gorilla/mux"
)

type TransactionHandler struct {
	service service.ITransactionService
}

func NewTransactionHandler(service service.ITransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	var (
		transaction models.Transaction
		resData     models.CreateTransactionResponse
	)

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ResponseBody{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	transactionID, err := h.service.CreateTransaction(ctx, transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ResponseBody{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	resData.ID = transactionID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.ResponseBody{
		Message: "Transaction created",
		Data:    resData,
	})
}

func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	var (
		transaction models.Transaction
	)

	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ResponseBody{
			Error:   true,
			Message: "ID must filled",
		})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ResponseBody{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	transaction.ID = id
	if err := h.service.UpdateTransaction(ctx, transaction); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ResponseBody{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.ResponseBody{
		Message: "Transaction Updated",
	})
}

func (h *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	var (
		transaction models.Transaction
	)

	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ResponseBody{
			Error:   true,
			Message: "ID must filled",
		})
		return
	}

	transaction, err := h.service.GetTransactionByID(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ResponseBody{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.ResponseBody{
		Message: "Success Get Transaction",
		Data:    transaction,
	})
}

func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ResponseBody{
			Error:   true,
			Message: "ID must filled",
		})
		return
	}

	if err := h.service.DeleteTransaction(ctx, id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ResponseBody{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.ResponseBody{
		Message: "Success Delete Transaction",
	})
}

func (h *TransactionHandler) GetAllTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	var (
		req models.GetAllTransactionRequest
	)

	query := r.URL.Query()
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		limit = 10
	}

	req.Limit = limit
	req.Page = page

	transactions, err := h.service.GetAllTransaction(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ResponseBody{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.ResponseBody{
		Message: "Success get all transaction",
		Data:    transactions,
	})
	return
}
