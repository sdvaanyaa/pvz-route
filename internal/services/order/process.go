package order

import (
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"time"
)

const returnWindow = 48 * time.Hour

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

	case "return":
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
	default:
		return ErrUnknownAction
	}

	if err = s.storage.UpdateOrder(order); err != nil {
		return err
	}

	return nil
}
