package storage

import (
	"errors"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
)

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderExpired       = errors.New("order expired")
)

type Storage interface {
	SaveOrder(order *models.Order) error
}
