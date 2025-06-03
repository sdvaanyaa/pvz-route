package order

import (
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"time"
)

const returnWindow = 48 * time.Hour

// Process executes an action ("issue" or "return") on an order
// for a given userID and orderID.
// It validates inputs, verifies order ownership and status,
// updates the order status accordingly, and saves changes to storage.
func (s *orderService) Process(userID, orderID, action string) error {
	if userID == "" {
		return ErrEmptyUserID
	}

	if orderID == "" {
		return ErrEmptyOrderID
	}

	order, err := s.storage.GetOrder(orderID)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return ErrOrderNotBelongsToUser
	}

	switch action {
	case "issue":
		if err = processIssue(order); err != nil {
			return err
		}
	case "return":
		if err = processReturn(order); err != nil {
			return err
		}
	default:
		return ErrUnknownAction
	}

	if err = s.storage.UpdateOrder(order); err != nil {
		return err
	}

	return nil
}

func processIssue(order *models.Order) error {
	if order.Status != models.StatusAccepted {
		return ErrOrderNotAccepted
	}

	if time.Now().After(order.StorageExpire) {
		return ErrStorageExpired
	}

	order.Status = models.StatusIssued
	now := time.Now()
	order.IssuedAt = &now
	order.History = append(order.History, models.OrderStatusChange{
		Status:    models.StatusIssued,
		Timestamp: now,
	})

	return nil
}

func processReturn(order *models.Order) error {
	if order.Status != models.StatusIssued || order.IssuedAt == nil {
		return ErrOrderNotIssued
	}

	if time.Since(*order.IssuedAt) > returnWindow {
		return ErrReturnPeriodExpired
	}

	order.Status = models.StatusReturned
	order.IssuedAt = nil
	now := time.Now()
	order.ReturnedAt = &now
	order.History = append(order.History, models.OrderStatusChange{
		Status:    models.StatusReturned,
		Timestamp: now,
	})

	return nil
}
