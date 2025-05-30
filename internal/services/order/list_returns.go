package order

import (
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"sort"
)

func (s *orderService) ListReturns(page, limit int) ([]*models.Order, error) {
	if page < 1 {
		return nil, ErrInvalidPageNumber
	}

	orders, err := s.storage.GetOrders()
	if err != nil {
		return nil, err
	}

	returns := filterReturnedOrders(orders)

	sort.Slice(returns, func(i, j int) bool {
		return returns[i].ReturnedAt.After(*returns[j].ReturnedAt)
	})

	returns = paginateOrders(returns, page, limit)

	return returns, nil
}
