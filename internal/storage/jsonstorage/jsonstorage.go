package jsonstorage

import (
	"encoding/json"
	"errors"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
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

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderExpired       = errors.New("order expired")
	ErrOrderNotFound      = errors.New("order not found")
	ErrStorageIO          = errors.New("failed to access storage")
)

// New creates a new Storage instance, initializing the storage directory
// and the orders JSON file if it does not exist.
func New(storagePath string) (*Storage, error) {
	if err := os.MkdirAll(storagePath, dirPerm); err != nil {
		return nil, err
	}

	ordersPath := filepath.Join(storagePath, "orders.json")

	if _, err := os.Stat(ordersPath); errors.Is(err, fs.ErrNotExist) {
		if err = os.WriteFile(ordersPath, []byte("[]"), filePerm); err != nil {
			return nil, err
		}
	}

	return &Storage{ordersPath}, nil
}

// SaveOrder saves a new order to the JSON storage.
// Validates that the order is not expired and does not already exist.
// Returns an error if validation fails or the operation cannot be completed.
func (s *Storage) SaveOrder(order *models.Order) error {
	if time.Now().After(order.StorageExpire) {
		return ErrOrderExpired
	}

	if _, err := s.GetOrder(order.ID); err == nil {
		return ErrOrderAlreadyExists
	} else if !errors.Is(err, ErrOrderNotFound) {
		return err
	}

	orders, err := s.GetOrders()
	if err != nil {
		return err
	}

	orders = append(orders, order)

	if err = s.SaveOrders(orders); err != nil {
		return err
	}

	return nil
}

// GetOrder retrieves an order by its ID from the JSON storage.
// Returns the order and nil if found, or nil and ErrOrderNotFound if not found.
func (s *Storage) GetOrder(orderID string) (*models.Order, error) {
	var orders []*models.Order
	if err := readJSON(s.ordersPath, &orders); err != nil {
		return nil, err
	}

	for _, o := range orders {
		if o.ID == orderID {
			return o, nil
		}
	}

	return nil, ErrOrderNotFound
}

// UpdateOrder updates an existing order in the JSON storage.
// Returns ErrOrderNotFound if the order does not exist or an error if the operation fails.
func (s *Storage) UpdateOrder(order *models.Order) error {
	orders, err := s.GetOrders()
	if err != nil {
		return err
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
		return ErrOrderNotFound
	}

	return s.SaveOrders(orders)
}

// GetOrdersByUser retrieves all orders for a given user from the JSON storage.
// Returns a slice of orders and an error if the operation fails.
func (s *Storage) GetOrdersByUser(userID string) ([]*models.Order, error) {
	var orders []*models.Order
	if err := readJSON(s.ordersPath, &orders); err != nil {
		return nil, err
	}

	var userOrders []*models.Order
	for _, o := range orders {
		if o.UserID == userID {
			userOrders = append(userOrders, o)
		}
	}

	return userOrders, nil
}

// SaveOrders saves a list of orders to the JSON storage.
// Overwrites the existing orders.json file.
// Returns an error if the operation fails.
func (s *Storage) SaveOrders(orders []*models.Order) error {
	if err := writeJSON(s.ordersPath, orders); err != nil {
		return err
	}

	return nil
}

// GetOrders retrieves all orders from the JSON storage.
// Returns a slice of orders and an error if the operation fails.
func (s *Storage) GetOrders() ([]*models.Order, error) {
	var orders []*models.Order
	if err := readJSON(s.ordersPath, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func readJSON(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return ErrStorageIO
	}

	if err = json.Unmarshal(data, v); err != nil {
		return ErrStorageIO
	}

	return nil
}

func writeJSON(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ErrStorageIO
	}

	if err = os.WriteFile(path, data, filePerm); err != nil {
		return ErrStorageIO
	}

	return nil
}
