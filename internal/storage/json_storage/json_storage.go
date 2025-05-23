package json_storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/storage"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	dirPerm  = 0755
	filePerm = 0644
)

type Storage struct {
	ordersPath string
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.json_storage.New"

	if err := os.MkdirAll(storagePath, dirPerm); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ordersPath := filepath.Join(storagePath, "orders.json")

	if _, err := os.Stat(ordersPath); errors.Is(err, fs.ErrNotExist) {
		if err = os.WriteFile(ordersPath, []byte("[]"), filePerm); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &Storage{ordersPath}, nil
}

func (s *Storage) SaveOrder(order *models.Order) error {
	const op = "storage.json_storage.SaveOrder"

	if time.Now().After(order.StorageExpire) {
		return fmt.Errorf("%s: %w", op, storage.ErrOrderExpired)
	}

	if _, err := s.GetOrder(order.ID); err == nil {
		return fmt.Errorf("%s: %w", op, storage.ErrOrderAlreadyExists)
	} else if !errors.Is(err, storage.ErrOrderNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	orders, err := s.GetOrders()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	orders = append(orders, order)

	if err = s.SaveOrders(orders); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetOrder(orderID string) (*models.Order, error) {
	const op = "storage.json_storage.GetOrder"

	var orders []*models.Order
	if err := readJSON(s.ordersPath, &orders); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for _, o := range orders {
		if o.ID == orderID {
			return o, nil
		}
	}

	return nil, fmt.Errorf("%s: %w", op, storage.ErrOrderNotFound)
}

func (s *Storage) UpdateOrder(order *models.Order) error {
	const op = "storage.json_storage.UpdateOrder"

	orders, err := s.GetOrders()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var updated bool
	for i, o := range orders {
		if order.ID == o.ID {
			orders[i] = order
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("%s: %w", op, storage.ErrOrderNotFound)
	}

	return s.SaveOrders(orders)
}

func (s *Storage) GetOrdersByUser(userID string) ([]*models.Order, error) {
	const op = "storage.json_storage.GetOrdersByUser"

	var orders []*models.Order
	if err := readJSON(s.ordersPath, &orders); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var userOrders []*models.Order
	for _, o := range orders {
		if o.UserID == userID {
			userOrders = append(userOrders, o)
		}
	}

	return userOrders, nil
}

func (s *Storage) SaveOrders(orders []*models.Order) error {
	const op = "storage.json_storage.SaveOrders"

	if err := writeJSON(s.ordersPath, orders); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetOrders() ([]*models.Order, error) {
	const op = "storage.json_storage.GetOrders"

	var orders []*models.Order
	if err := readJSON(s.ordersPath, &orders); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return orders, nil
}

func readJSON(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	return nil
}

func writeJSON(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	if err := os.WriteFile(path, data, filePerm); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	return nil
}
