package order

import (
	"encoding/json"
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
	if path == "" {
		return 0, ErrEmptyFilePath
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	var orders []importOrder

	if err = json.Unmarshal(data, &orders); err != nil {
		return 0, err
	}

	if len(orders) == 0 {
		return 0, ErrEmptyImportFile
	}

	existOrders, err := s.storage.GetOrders()
	if err != nil {
		return 0, err
	}

	existIDs := make(map[string]struct{}, len(existOrders))
	for _, order := range existOrders {
		existIDs[order.ID] = struct{}{}
	}

	newOrders := make([]*models.Order, 0, len(orders))
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

		newOrders = append(newOrders, newOrder)
		existIDs[order.ID] = struct{}{}
	}

	if len(newOrders) == 0 {
		return 0, ErrEmptyValidOrders
	}

	importCount := len(newOrders)

	allOrders := make([]*models.Order, 0, len(newOrders)+len(existOrders))
	allOrders = append(allOrders, newOrders...)
	allOrders = append(allOrders, existOrders...)

	if err = s.storage.SaveOrders(allOrders); err != nil {
		return 0, err
	}

	return importCount, nil
}

func isOrderValid(order importOrder) bool {
	return order.ID != "" &&
		order.UserID != "" &&
		order.StorageExpire != nil &&
		order.StorageExpire.After(time.Now())
}
