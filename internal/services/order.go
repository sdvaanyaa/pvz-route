package services

import (
	"errors"
	"fmt"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/storage"
	"sort"
	"time"
)

var (
	ErrEmptyOrderID          = errors.New("order ID must not be empty")
	ErrEmptyUserID           = errors.New("user ID must not be empty")
	ErrInvalidDeadlineFormat = errors.New("deadline must be in YYYY-MM-DD format")
	ErrStorageNotExpired     = errors.New("cannot return order: storage period has not expired yet")
	ErrStorageExpired        = errors.New("cannot issue order: storage period expired")
	ErrOrderIssued           = errors.New("order has already been issued")
	ErrOrderNotIssued        = errors.New("order has not yet been issued")
	ErrUnknownAction         = errors.New("action must be specified: 'issue' or 'return'")
	ErrOrderNotBelongsToUser = errors.New("order does not belong to user")
	ErrOrderNotAccepted      = errors.New("order has not been accepted")
	ErrReturnPeriodExpired   = errors.New("return period exceeded: more than 48 hours since issue")
)

type OrderService struct {
	storage storage.Storage
}

func New(storage storage.Storage) *OrderService {
	return &OrderService{storage: storage}
}

func (s *OrderService) AcceptOrder(orderID, userID, expire string) error {
	const op = "services.order.AcceptOrder"

	if orderID == "" {
		return fmt.Errorf("%s: %w", op, ErrEmptyOrderID)
	}

	if userID == "" {
		return fmt.Errorf("%s: %w", op, ErrEmptyUserID)
	}

	storageExpire, err := time.Parse(time.DateOnly, expire)
	if err != nil {
		return fmt.Errorf("%s: %w", op, ErrInvalidDeadlineFormat)
	}

	order := &models.Order{
		ID:            orderID,
		UserID:        userID,
		StorageExpire: storageExpire,
		Status:        models.StatusAccepted,
		CreatedAt:     time.Now(),
	}

	return s.storage.SaveOrder(order)
}

func (s *OrderService) ReturnOrder(orderID string) error {
	const op = "services.order.ReturnOrder"

	if orderID == "" {
		return fmt.Errorf("%s: %w", op, ErrEmptyOrderID)
	}

	order, err := s.storage.GetOrder(orderID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if order.StorageExpire.After(time.Now()) {
		return fmt.Errorf("%s: %w", op, ErrStorageNotExpired)
	}

	if order.Status == models.StatusIssued && order.IssuedAt != nil {
		return fmt.Errorf("%s: %w", op, ErrOrderIssued)
	}

	return s.storage.DeleteOrder(orderID)
}

func (s *OrderService) ProcessOrder(userID, orderID, action string) error {
	const op = "services.order.ProcessOrder"

	if userID == "" {
		return fmt.Errorf("%s: %w", op, ErrEmptyUserID)
	}

	if orderID == "" {
		return fmt.Errorf("%s: %w", op, ErrEmptyOrderID)
	}

	order, err := s.storage.GetOrder(orderID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if order.UserID != userID {
		return fmt.Errorf("%s: %w", op, ErrOrderNotBelongsToUser)
	}

	switch action {
	case "issue":
		if order.Status != models.StatusAccepted {
			return fmt.Errorf("%s: %w", op, ErrOrderNotAccepted)
		}

		if time.Now().After(order.StorageExpire) {
			return fmt.Errorf("%s: %w", op, ErrStorageExpired)
		}

		order.Status = models.StatusIssued
		now := time.Now()
		order.IssuedAt = &now

		if err = s.storage.UpdateOrder(order); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	case "return":
		if order.Status != models.StatusIssued || order.IssuedAt == nil {
			return fmt.Errorf("%s: %w", op, ErrOrderNotIssued)
		}

		if time.Since(*order.IssuedAt).Hours() > 48 {
			return fmt.Errorf("%s: %w", op, ErrReturnPeriodExpired)
		}

		order.Status = models.StatusReturned
		order.IssuedAt = nil
		if err = s.storage.UpdateOrder(order); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	default:
		return fmt.Errorf("%s: %w", op, ErrUnknownAction)
	}

	return nil
}

func (s *OrderService) ListOrders(userID string, inPVZ bool, last, page, limit int) ([]*models.Order, int, error) {
	const op = "services.order.ListOrders"

	if userID == "" {
		return nil, 0, fmt.Errorf("%s: %w", op, ErrEmptyUserID)
	}

	orders, err := s.storage.GetOrdersByUser(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	if inPVZ {
		filtered := make([]*models.Order, 0, len(orders))

		for _, order := range orders {
			if order.Status == models.StatusAccepted {
				filtered = append(filtered, order)
			}
		}

		orders = filtered
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
