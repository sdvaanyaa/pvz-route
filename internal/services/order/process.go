package order

import (
	"fmt"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"time"
)

func (s *orderService) Process(userID, orderID, action string) error {
	const op = "services.order.Process"

	if userID == "" {
		return fmt.Errorf("%s: %w", op, ErrEmptyUserID)
	}

	if orderID == "" {
		return fmt.Errorf("%s: %w", op, ErrEmptyOrderID)
	}

	order, err := s.storage.GetOrder(orderID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if order.UserID != userID {
		return fmt.Errorf("%s: %w", op, ErrOrderNotBelongsToUser)
	}

	switch action {
	case "issue":
		if order.Status != models.StatusAccepted {
			return fmt.Errorf("%s: %w", op, ErrOrderNotAccepted)
		}

		if time.Now().After(order.StorageExpire) {
			return fmt.Errorf("%s: %w", op, ErrStorageExpired)
		}

		order.Status = models.StatusIssued
		now := time.Now()
		order.IssuedAt = &now
		order.History = append(order.History, models.OrderStatusChange{
			Status:    models.StatusIssued,
			Timestamp: now,
		})

	case "return":
		if order.Status != models.StatusIssued || order.IssuedAt == nil {
			return fmt.Errorf("%s: %w", op, ErrOrderNotIssued)
		}

		if time.Since(*order.IssuedAt).Hours() > 48 {
			return fmt.Errorf("%s: %w", op, ErrReturnPeriodExpired)
		}

		order.Status = models.StatusReturned
		order.IssuedAt = nil
		now := time.Now()
		order.ReturnedAt = &now
		order.History = append(order.History, models.OrderStatusChange{
			Status:    models.StatusReturned,
			Timestamp: now,
		})
	default:
		return fmt.Errorf("%s: %w", op, ErrUnknownAction)
	}

	if err = s.storage.UpdateOrder(order); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
