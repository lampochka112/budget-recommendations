package handlers

import (
	"budget-recommendations/internal"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TransactionHandler struct {
	service internal.TransactionService
}

func NewTransactionHandler(service internal.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

type CreateTransactionRequest struct {
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
}

type TransactionResponse struct {
	ID          uuid.UUID `json:"id"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	UserID      uuid.UUID `json:"user_id"`
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Валидация
	if req.Amount <= 0 {
		http.Error(w, "Amount must be positive", http.StatusBadRequest)
		return
	}
	if req.Category == "" {
		http.Error(w, "Category is required", http.StatusBadRequest)
		return
	}
	if req.Type != "income" && req.Type != "expense" {
		http.Error(w, "Type must be 'income' or 'expense'", http.StatusBadRequest)
		return
	}

	// TODO: Получать userID из JWT токена (пока заглушка)
	userID := uuid.MustParse("a5c5eaa5-1234-5678-9012-1ea3d5bb1234")

	transaction := &internal.Transaction{
		ID:          uuid.New(),
		Amount:      req.Amount,
		Category:    req.Category,
		Date:        req.Date,
		Description: req.Description,
		Type:        req.Type,
		UserID:      userID,
	}

	if err := h.service.CreateTransaction(r.Context(), transaction); err != nil {
		http.Error(w, "Failed to create transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(TransactionResponse{
		ID:          transaction.ID,
		Amount:      transaction.Amount,
		Category:    transaction.Category,
		Date:        transaction.Date,
		Description: transaction.Description,
		Type:        transaction.Type,
		UserID:      transaction.UserID,
	})
}

func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	// TODO: Получать userID из JWT токена
	userID := uuid.MustParse("a5c5eaa5-1234-5678-9012-1ea3d5bb1234")

	query := r.URL.Query()
	filters := make(map[string]interface{})

	// Парсим параметры запроса
	if category := query.Get("category"); category != "" {
		filters["category"] = category
	}
	if transactionType := query.Get("type"); transactionType != "" {
		filters["type"] = transactionType
	}
	if startDate := query.Get("start_date"); startDate != "" {
		if date, err := time.Parse("2006-01-02", startDate); err == nil {
			filters["start_date"] = date
		}
	}
	if endDate := query.Get("end_date"); endDate != "" {
		if date, err := time.Parse("2006-01-02", endDate); err == nil {
			filters["end_date"] = date
		}
	}

	// Пагинация
	if limit := query.Get("limit"); limit != "" {
		if limitInt, err := strconv.Atoi(limit); err == nil && limitInt > 0 {
			filters["limit"] = limitInt
		}
	}
	if page := query.Get("page"); page != "" {
		if pageInt, err := strconv.Atoi(page); err == nil && pageInt > 0 {
			if limit, ok := filters["limit"].(int); ok {
				filters["offset"] = (pageInt - 1) * limit
			}
		}
	}

	transactions, err := h.service.GetTransactions(r.Context(), userID, filters)
	if err != nil {
		http.Error(w, "Failed to get transactions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]TransactionResponse, len(transactions))
	for i, t := range transactions {
		response[i] = TransactionResponse{
			ID:          t.ID,
			Amount:      t.Amount,
			Category:    t.Category,
			Date:        t.Date,
			Description: t.Description,
			Type:        t.Type,
			UserID:      t.UserID,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Получать userID из JWT токена
	userID := uuid.MustParse("a5c5eaa5-1234-5678-9012-1ea3d5bb1234")

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.GetTransactionByID(r.Context(), id, userID)
	if err != nil {
		http.Error(w, "Transaction not found: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TransactionResponse{
		ID:          transaction.ID,
		Amount:      transaction.Amount,
		Category:    transaction.Category,
		Date:        transaction.Date,
		Description: transaction.Description,
		Type:        transaction.Type,
		UserID:      transaction.UserID,
	})
}

func (h *TransactionHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}