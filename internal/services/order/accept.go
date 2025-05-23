package order

import (
	"fmt"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"time"
)

func (s *Service) AcceptOrder(orderID, userID, expire string) error {
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
