package order

import (
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/storage"
)

// Service defines the interface for order management operations,
// including accepting, processing, listing, importing, and returning orders.
type Service interface {
	// Accept creates and accepts a new order with the specified parameters.
	Accept(orderID, userID, expire string, weight, price float64, packageType string) (*models.Order, error)

	// History returns all status changes of all orders.
	History() ([]*HistoryEntry, error)

	// ListOrders lists orders for a user, optionally filtered and paginated.
	ListOrders(userID string, inPVZ bool, last, page, limit int) ([]*models.Order, int, error)

	// ImportOrders imports orders from a JSON file at the given path.
	ImportOrders(path string) (int, error)

	// ListReturns returns a list of returned orders with optional pagination support.
	ListReturns(page, limit int) ([]*models.Order, error)

	// Process executes an action ("issue" or "return") on a specified order.
	Process(userID, orderID, action string) error

	// Return updates the order status to archive when the storage period has ended.
	Return(orderID string) error

	// Scroll fetches active orders for a user with support for infinite scrolling.
	Scroll(userID, lastID string, limit int) ([]*models.Order, string, error)
}

type orderService struct {
	storage storage.Storage
}

// New creates a new order service with the provided storage implementation.
func New(storage storage.Storage) Service {
	return &orderService{storage: storage}
}
