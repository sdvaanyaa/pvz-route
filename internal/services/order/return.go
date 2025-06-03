package order

import (
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"time"
)

// Return marks an order as archived if its storage time has ended,
// and it is not currently issued. It updates the order status and saves changes.
func (s *orderService) Return(orderID string) error {
	if orderID == "" {
		return ErrEmptyOrderID
	}

	order, err := s.storage.GetOrder(orderID)
	if err != nil {
		return err
	}

	if order.StorageExpire.After(time.Now()) {
		return ErrStorageNotExpired
	}

	if order.Status == models.StatusIssued && order.IssuedAt != nil {
		return ErrOrderIssued
	}

	now := time.Now()
	order.Status = models.StatusArchived
	order.ArchivedAt = &now
	order.History = append(order.History, models.OrderStatusChange{
		Status:    models.StatusArchived,
		Timestamp: now,
	})

	if err = s.storage.UpdateOrder(order); err != nil {
		return err
	}

	return nil
}
