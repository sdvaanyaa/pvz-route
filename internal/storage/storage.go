package storage

import (
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
)

// Storage defines methods for managing orders in a storage.
type Storage interface {
	// SaveOrder saves a new order to the storage.
	// Returns an error if the operation fails.
	SaveOrder(order *models.Order) error

	// GetOrder retrieves an order by its ID.
	// Returns the order and nil if found, or nil and an error if not found.
	GetOrder(orderID string) (*models.Order, error)

	// UpdateOrder updates an existing order in the storage.
	// Returns an error if the order is not found or the operation fails.
	UpdateOrder(order *models.Order) error

	// GetOrdersByUser retrieves all orders for a given user.
	// Returns a slice of orders and an error if the operation fails.
	GetOrdersByUser(userID string) ([]*models.Order, error)

	// SaveOrders saves a list of orders to the storage.
	// Returns an error if the operation fails.
	SaveOrders(orders []*models.Order) error

	// GetOrders retrieves all orders from the storage.
	// Returns a slice of orders and an error if the operation fails.
	GetOrders() ([]*models.Order, error)
}
