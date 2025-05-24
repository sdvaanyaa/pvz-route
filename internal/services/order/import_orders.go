package order

import (
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"os"
	"time"
)

type importOrder struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	StorageExpire *time.Time `json:"storage_deadline"`
}

func (s *orderService) ImportOrders(path string) (int, error) {
	const op = "services.order.ImportOrders"

	if path == "" {
		return 0, fmt.Errorf("%s: %w", op, ErrEmptyFilePath)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var orders []importOrder

	if err = json.Unmarshal(data, &orders); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if len(orders) == 0 {
		return 0, fmt.Errorf("%s: %w", op, ErrEmptyImportFile)
	}

	existOrders, err := s.storage.GetOrders()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	existIDs := make(map[string]struct{}, len(existOrders))
	for _, order := range existOrders {
		existIDs[order.ID] = struct{}{}
	}

	validOrders := make([]*models.Order, 0, len(orders))
	now := time.Now()

	for _, order := range orders {
		if _, ok := existIDs[order.ID]; ok {
			continue
		}

		if !isOrderValid(order) {
			continue
		}

		newOrder := &models.Order{
			ID:            order.ID,
			UserID:        order.UserID,
			StorageExpire: *order.StorageExpire,
			Status:        models.StatusAccepted,
			CreatedAt:     now,
			History: []models.OrderStatusChange{
				{
					Status:    models.StatusAccepted,
					Timestamp: now,
				},
			},
		}

		validOrders = append(validOrders, newOrder)
		existIDs[order.ID] = struct{}{}
	}

	if len(validOrders) == 0 {
		return 0, fmt.Errorf("%s: %w", op, ErrEmptyValidOrders)
	}

	if err = s.storage.SaveOrders(validOrders); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return len(validOrders), nil
}

func isOrderValid(order importOrder) bool {
	return order.ID != "" &&
		order.UserID != "" &&
		order.StorageExpire != nil &&
		order.StorageExpire.After(time.Now())
}
