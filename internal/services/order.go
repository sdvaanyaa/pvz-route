package services

import (
	"errors"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/storage"
	"time"
)

var (
	ErrEmptyOrderID          = errors.New("empty orderID")
	ErrEmptyRecipientID      = errors.New("empty recipientID")
	ErrInvalidDeadlineFormat = errors.New("invalid deadline format")
)

type OrderService struct {
	storage storage.Storage
}

func New(storage storage.Storage) *OrderService {
	return &OrderService{storage: storage}
}

func (s *OrderService) AcceptOrder(orderID, userID, expire string) error {
	if orderID == "" {
		return ErrEmptyOrderID
	}

	if userID == "" {
		return ErrEmptyRecipientID
	}

	storageExpire, err := time.Parse(time.DateOnly, expire)
	if err != nil {
		return ErrInvalidDeadlineFormat
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
