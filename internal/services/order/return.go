package order

import (
	"fmt"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"time"
)

func (s *Service) ReturnOrder(orderID string) error {
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
