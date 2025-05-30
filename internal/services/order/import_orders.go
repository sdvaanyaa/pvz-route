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

	orders, err := readOrdersFromFile(path)
	if err != nil {
		return 0, err
	}

	if len(orders) == 0 {
		return 0, ErrEmptyImportFile
	}

	existOrders, err := s.storage.GetOrders()
	if err != nil {
		return 0, err
	}

	newOrders := prepareNewOrders(orders, existOrders)

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

func prepareNewOrders(imports []*importOrder, existOrders []*models.Order) []*models.Order {
	existIDs := make(map[string]struct{}, len(existOrders))
	for _, o := range existOrders {
		existIDs[o.ID] = struct{}{}
	}

	now := time.Now()
	newOrders := make([]*models.Order, 0, len(imports))

	for _, o := range imports {
		if _, exists := existIDs[o.ID]; exists || !o.IsValid() {
			continue
		}
		newOrders = append(newOrders, convertToModelOrder(o, now))
		existIDs[o.ID] = struct{}{}
	}

	return newOrders
}

func convertToModelOrder(o *importOrder, now time.Time) *models.Order {
	return &models.Order{
		ID:            o.ID,
		UserID:        o.UserID,
		StorageExpire: *o.StorageExpire,
		Status:        models.StatusAccepted,
		CreatedAt:     now,
		History: []models.OrderStatusChange{
			{
				Status:    models.StatusAccepted,
				Timestamp: now,
			},
		},
	}
}

func readOrdersFromFile(path string) ([]*importOrder, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var orders []*importOrder

	if err = json.Unmarshal(data, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *importOrder) IsValid() bool {
	return o.ID != "" &&
		o.UserID != "" &&
		o.StorageExpire != nil &&
		o.StorageExpire.After(time.Now())
}
