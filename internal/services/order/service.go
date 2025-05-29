package order

import (
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/storage"
)

type Service interface {
	Accept(orderID, userID, expire string, weight, price float64, packageType string) (*models.Order, error)
	History() ([]*HistoryEntry, error)
	ListOrders(userID string, inPVZ bool, last, page, limit int) ([]*models.Order, int, error)
	ImportOrders(path string) (int, error)
	ListReturns(page, limit int) ([]*models.Order, error)
	Process(userID, orderID, action string) error
	Return(orderID string) error
	Scroll(userID, lastID string, limit int) ([]*models.Order, string, error)
}

type orderService struct {
	storage storage.Storage
}

func New(storage storage.Storage) Service {
	return &orderService{storage: storage}
}
