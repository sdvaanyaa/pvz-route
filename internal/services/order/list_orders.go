package order

import (
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"sort"
)

func (s *orderService) ListOrders(userID string, inPVZ bool, last, page, limit int) ([]*models.Order, int, error) {
	if userID == "" {
		return nil, 0, ErrEmptyUserID
	}

	if page < 1 {
		return nil, 0, ErrInvalidPageNumber
	}

	orders, err := s.storage.GetOrdersByUser(userID)
	if err != nil {
		return nil, 0, err
	}

	orders = filterActiveOrders(orders)

	if inPVZ {
		orders = filterAcceptedOrders(orders)
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].CreatedAt.After(orders[j].CreatedAt)
	})

	if last > 0 {
		orders = orders[:min(last, len(orders))]
	}

	total := len(orders)

	orders = paginateOrders(orders, page, limit)

	return orders, total, nil
}
