package service

import (
	"budget-recommendations/internal"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type TransactionService struct {
	storage internal.TransactionStorage
}

func NewTransactionService(storage internal.TransactionStorage) *TransactionService {
	return &TransactionService{storage: storage}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, transaction *internal.Transaction) error {
	// Валидация
	if transaction.Amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if transaction.Category == "" {
		return fmt.Errorf("category is required")
	}
	if transaction.Type != "income" && transaction.Type != "expense" {
		return fmt.Errorf("type must be 'income' or 'expense'")
	}

	return s.storage.CreateTransaction(ctx, transaction)
}

func (s *TransactionService) GetTransactions(ctx context.Context, userID uuid.UUID, filters map[string]interface{}) ([]internal.Transaction, error) {
	// Добавляем userID в фильтры для безопасности
	if filters == nil {
		filters = make(map[string]interface{})
	}
	
	return s.storage.GetTransactions(ctx, userID, filters)
}

func (s *TransactionService) GetTransactionByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*internal.Transaction, error) {
	transaction, err := s.storage.GetTransactionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверяем, что транзакция принадлежит пользователю
	if transaction.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	return transaction, nil
}

func (s *TransactionService) UpdateTransaction(ctx context.Context, transaction *internal.Transaction, userID uuid.UUID) error {
	// Проверяем, что транзакция принадлежит пользователю
	existing, err := s.GetTransactionByID(ctx, transaction.ID, userID)
	if err != nil {
		return err
	}

	// Обновляем только разрешенные поля
	existing.Amount = transaction.Amount
	existing.Category = transaction.Category
	existing.Date = transaction.Date
	existing.Description = transaction.Description
	existing.Type = transaction.Type

	return s.storage.UpdateTransaction(ctx, existing)
}

func (s *TransactionService) DeleteTransaction(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Проверяем существование транзакции и права доступа
	_, err := s.GetTransactionByID(ctx, id, userID)
	if err != nil {
		return err
	}

	return s.storage.DeleteTransaction(ctx, id, userID)
}