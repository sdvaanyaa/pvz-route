package order

import (
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"time"
)

func (s *orderService) Accept(orderID, userID, expire string) error {
	if orderID == "" {
		return ErrEmptyOrderID
	}

	if userID == "" {
		return ErrEmptyUserID
	}

	storageExpire, err := time.Parse(time.DateOnly, expire)
	if err != nil {
		return ErrInvalidDeadlineFormat
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
