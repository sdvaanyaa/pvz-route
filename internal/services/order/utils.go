package order

import "gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"

func paginateOrders(orders []*models.Order, page, limit int) []*models.Order {
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
	return orders
}

func filterActiveOrders(orders []*models.Order) []*models.Order {
	active := make([]*models.Order, 0, len(orders))
	for _, order := range orders {
		if order.Status != models.StatusArchived {
			active = append(active, order)
		}
	}
	return active
}

func filterAcceptedOrders(orders []*models.Order) []*models.Order {
	accepted := make([]*models.Order, 0, len(orders))
	for _, order := range orders {
		if order.Status == models.StatusAccepted {
			accepted = append(accepted, order)
		}
	}
	return accepted
}

func filterReturnedOrders(orders []*models.Order) []*models.Order {
	returned := make([]*models.Order, 0, len(orders))
	for _, order := range orders {
		if order.Status == models.StatusReturned && order.ReturnedAt != nil {
			returned = append(returned, order)
		}
	}
	return returned
}
