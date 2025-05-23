package order

import (
	"fmt"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"time"
)

func (s *Service) Accept(orderID, userID, expire string) error {
	const op = "services.order.Accept"

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

	now := time.Now()
	order := &models.Order{
		ID:            orderID,
		UserID:        userID,
		StorageExpire: storageExpire,
		Status:        models.StatusAccepted,
		CreatedAt:     now,
		History: []models.OrderStatusChange{
			{
				Status:    models.StatusAccepted,
				Timestamp: now,
			},
		},
	}

	return s.storage.SaveOrder(order)
}
