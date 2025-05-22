package storage

import (
	"errors"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
)

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderExpired       = errors.New("order expired")
	ErrOrderNotFound      = errors.New("order not found")
)

type Storage interface {
	SaveOrder(order *models.Order) error
	GetOrder(orderID string) (*models.Order, error)
	DeleteOrder(orderID string) error
	UpdateOrder(order *models.Order) error

	GetOrdersByUser(userID string) ([]*models.Order, error)

	SaveOrders(orders []*models.Order) error
	GetOrders() ([]*models.Order, error)
}
