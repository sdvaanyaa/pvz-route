package order

import (
	"fmt"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/packaging"
	"time"
)

func (s *orderService) Accept(
	orderID, userID, expire string,
	weight, price float64,
	packageType string,
) (*models.Order, error) {
	if err := validateInputs(orderID, userID, price, weight); err != nil {
		return nil, err
	}

	storageExpire, err := time.Parse(time.DateOnly, expire)
	if err != nil {
		return nil, ErrInvalidDeadlineFormat
	}

	if packageType == "" {
		packageType = packaging.PackageNone
	}

	finalPrice, err := prepareAndValidatePackage(packageType, price, weight)
	if err != nil {
		return nil, err
	}

	order := newOrder(
		orderID,
		userID,
		storageExpire,
		packageType,
		weight,
		finalPrice,
	)

	return order, s.storage.SaveOrder(order)
}

func newOrder(
	orderID, userID string,
	storageExpire time.Time,
	packageType string,
	weight, price float64,
) *models.Order {
	now := time.Now()
	return &models.Order{
		ID:            orderID,
		UserID:        userID,
		StorageExpire: storageExpire,
		Status:        models.StatusAccepted,
		CreatedAt:     now,
		History: []models.OrderStatusChange{
			{
				Status:    models.StatusAccepted,
				Timestamp: now,
			},
		},
		Weight:      weight,
		Price:       price,
		PackageType: packageType,
	}
}

func validateInputs(orderID, userID string, price, weight float64) error {
	if orderID == "" {
		return ErrEmptyOrderID
	}

	if userID == "" {
		return ErrEmptyUserID
	}

	if price <= 0 {
		return fmt.Errorf("%w: got price=%.2f", ErrInvalidPrice, price)
	}

	if weight <= 0 {
		return fmt.Errorf("%w: got weight=%.2f", ErrInvalidWeight, weight)
	}

	return nil
}

func prepareAndValidatePackage(packageType string, price, weight float64) (float64, error) {
	strategy, err := packaging.GetPackageStrategy(packageType)
	if err != nil {
		return 0, err
	}

	if err = strategy.ValidateWeight(weight); err != nil {
		return 0, err
	}

	finalPrice := strategy.CalculatePrice(price)

	return finalPrice, nil
}
