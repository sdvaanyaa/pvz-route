package storage

import (
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
)

type Storage interface {
	SaveOrder(order *models.Order) error
	GetOrder(orderID string) (*models.Order, error)
	UpdateOrder(order *models.Order) error

	GetOrdersByUser(userID string) ([]*models.Order, error)

	SaveOrders(orders []*models.Order) error
	GetOrders() ([]*models.Order, error)
}
