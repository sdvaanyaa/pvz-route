package order

import (
	"fmt"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"sort"
)

func (s *Service) ListOrders(userID string, inPVZ bool, last, page, limit int) ([]*models.Order, int, error) {
	const op = "services.order.ListOrders"

	if userID == "" {
		return nil, 0, fmt.Errorf("%s: %w", op, ErrEmptyUserID)
	}

	orders, err := s.storage.GetOrdersByUser(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	activeOrders := make([]*models.Order, 0, len(orders))
	for _, order := range orders {
		if order.Status != models.StatusArchived {
			activeOrders = append(activeOrders, order)
		}
	}

	orders = activeOrders

	if inPVZ {
		inPVZOrders := make([]*models.Order, 0, len(orders))

		for _, order := range orders {
			if order.Status == models.StatusAccepted {
				inPVZOrders = append(inPVZOrders, order)
			}
		}

		orders = inPVZOrders
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].CreatedAt.After(orders[j].CreatedAt)
	})

	if last > 0 {
		orders = orders[:min(last, len(orders))]
	}

	total := len(orders)

	if limit > 0 {
		start := (page - 1) * limit
		if start >= len(orders) {
			orders = []*models.Order{}
		} else {
			end := start + limit
			if end > len(orders) {
				end = len(orders)
			}
			orders = orders[start:end]
		}
	}

	return orders, total, nil
}
