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
	ordersPath  string
	returnsPath string
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.json_storage.New"

	if err := os.MkdirAll(storagePath, dirPerm); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ordersPath := filepath.Join(storagePath, "orders.json")
	returnsPath := filepath.Join(storagePath, "returns.json")

	for _, path := range []string{ordersPath, returnsPath} {
		if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
			if err = os.WriteFile(path, []byte("[]"), filePerm); err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
		}
	}

	return &Storage{ordersPath, returnsPath}, nil
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

func (s *Storage) DeleteOrder(orderID string) error {
	const op = "storage.json_storage.DeleteOrder"

	orders, err := s.GetOrders()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for i, o := range orders {
		if o.ID == orderID {
			orders = append(orders[:i], orders[i+1:]...)
			return s.SaveOrders(orders)
		}
	}

	return fmt.Errorf("%s: %w", op, storage.ErrOrderNotFound)
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
