package json_storage

import (
	"encoding/json"
	"errors"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/storage"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/utils/errwrap"
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

func New(storagePath string) (db *Storage, err error) {
	defer func() { err = errwrap.WrapIfErr("failed to initialize JSON storage", err) }()

	if err := os.MkdirAll(storagePath, dirPerm); err != nil {
		return nil, err
	}

	ordersPath := filepath.Join(storagePath, "orders.json")
	returnsPath := filepath.Join(storagePath, "returns.json")

	for _, path := range []string{ordersPath, returnsPath} {
		if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
			if err := os.WriteFile(path, []byte("[]"), filePerm); err != nil {
				return nil, err
			}
		}
	}

	return &Storage{ordersPath, returnsPath}, nil
}

func (s *Storage) SaveOrder(order *models.Order) (err error) {
	defer func() { err = errwrap.WrapIfErr("failed to save order", err) }()

	if time.Now().After(order.StorageExpire) {
		return storage.ErrOrderExpired
	}

	var orders []models.Order
	if err = readJSON(s.ordersPath, &orders); err != nil {
		return err
	}

	for _, o := range orders {
		if o.ID == order.ID {
			return storage.ErrOrderAlreadyExists
		}
	}

	orders = append(orders, *order)

	if err := writeJSON(s.ordersPath, orders); err != nil {
		return err
	}

	return nil
}

func readJSON(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errwrap.Wrap("failed to read from JSON", err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return errwrap.Wrap("failed to unmarshal JSON", err)
	}

	return nil
}

func writeJSON(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return errwrap.Wrap("failed to marshal JSON", err)
	}

	if err := os.WriteFile(path, data, filePerm); err != nil {
		return errwrap.Wrap("failed to write to JSON", err)
	}

	return nil
}
