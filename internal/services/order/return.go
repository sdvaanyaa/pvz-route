package order

import (
	"fmt"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"time"
)

func (s *orderService) Return(orderID string) error {
	const op = "services.order.Return"

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

	now := time.Now()
	order.Status = models.StatusArchived
	order.ArchivedAt = &now
	order.History = append(order.History, models.OrderStatusChange{
		Status:    models.StatusArchived,
		Timestamp: now,
	})

	if err = s.storage.UpdateOrder(order); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
