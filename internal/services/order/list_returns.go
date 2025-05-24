package order

import (
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"sort"
)

func (s *orderService) ListReturns(page, limit int) ([]*models.Order, error) {
	if page < 1 {
		return nil, ErrInvalidPageNumber
	}

	if limit < 1 {
		return nil, ErrInvalidLastNumber
	}

	orders, err := s.storage.GetOrders()
	if err != nil {
		return nil, err
	}

	returns := make([]*models.Order, 0, len(orders))
	for _, order := range orders {
		if order.Status == models.StatusReturned && order.ReturnedAt != nil {
			returns = append(returns, order)
		}
	}

	sort.Slice(returns, func(i, j int) bool {
		return returns[i].ReturnedAt.After(*returns[j].ReturnedAt)
	})

	start := (page - 1) * limit
	if start >= len(orders) {
		returns = []*models.Order{}
	} else {
		end := start + limit
		if end > len(returns) {
			end = len(returns)
		}
		returns = returns[start:end]
	}

	return returns, nil
}
