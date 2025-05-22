package services

import (
	"errors"
	"fmt"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/storage"
	"time"
)

var (
	ErrEmptyOrderID          = errors.New("empty orderID")
	ErrEmptyUserID           = errors.New("empty userID")
	ErrInvalidDeadlineFormat = errors.New("invalid deadline format")
	ErrStorageNotExpired     = errors.New("storage period has not yet expired")
	ErrOrderIssued           = errors.New("order issued")
)

type OrderService struct {
	storage storage.Storage
}

func New(storage storage.Storage) *OrderService {
	return &OrderService{storage: storage}
}

func (s *OrderService) AcceptOrder(orderID, userID, expire string) error {
	const op = "services.order.AcceptOrder"

	if orderID == "" {
		return fmt.Errorf("%s: %w", op, ErrEmptyOrderID)
	}

	if userID == "" {
		return fmt.Errorf("%s: %w", op, ErrEmptyUserID)
	}

	storageExpire, err := time.Parse(time.DateOnly, expire)
	if err != nil {
		return fmt.Errorf("%s: %w", op, ErrInvalidDeadlineFormat)
	}

	order := &models.Order{
		ID:            orderID,
		UserID:        userID,
		StorageExpire: storageExpire,
		Status:        models.StatusAccepted,
		CreatedAt:     time.Now(),
	}

	return s.storage.SaveOrder(order)
}

func (s *OrderService) ReturnOrder(orderID string) error {
	const op = "services.order.ReturnOrder"

	if orderID == "" {
		return fmt.Errorf("%s: %w", op, ErrEmptyOrderID)
	}

	order, err := s.storage.GetOrder(orderID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if order.StorageExpire.After(time.Now()) {
		return fmt.Errorf("%s: %w", op, ErrStorageNotExpired)
	}

	if order.Status == models.StatusIssued && order.IssuedAt != nil {
		return fmt.Errorf("%s: %w", op, ErrOrderIssued)
	}

	return s.storage.DeleteOrder(orderID)
}
