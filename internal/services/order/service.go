package order

import (
	"errors"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/storage"
)

var (
	ErrEmptyOrderID          = errors.New("order ID must not be empty")
	ErrEmptyUserID           = errors.New("user ID must not be empty")
	ErrInvalidPageNumber     = errors.New("page number must be 1 or greater")
	ErrInvalidLastNumber     = errors.New("last number must be 1 or greater")
	ErrInvalidDeadlineFormat = errors.New("deadline must be in YYYY-MM-DD format")
	ErrStorageNotExpired     = errors.New("cannot return order: storage period has not expired yet")
	ErrStorageExpired        = errors.New("cannot issue order: storage period expired")
	ErrOrderIssued           = errors.New("order has already been issued")
	ErrOrderNotIssued        = errors.New("order has not yet been issued")
	ErrUnknownAction         = errors.New("action must be specified: 'issue' or 'return'")
	ErrOrderNotBelongsToUser = errors.New("order does not belong to user")
	ErrOrderNotAccepted      = errors.New("order has not been accepted")
	ErrReturnPeriodExpired   = errors.New("return period exceeded: more than 48 hours since issue")
	ErrEmptyFilePath         = errors.New("file path must not be empty")
	ErrEmptyImportFile       = errors.New("import file is empty")
	ErrEmptyValidOrders      = errors.New("import file does not contain valid orders")
	ErrInvalidLimitNumber    = errors.New("limit number must be 1 or greater")
)

type Service interface {
	Accept(orderID, userID, expire string) error
	History() ([]*HistoryEntry, error)
	ListOrders(userID string, inPVZ bool, last, page, limit int) ([]*models.Order, int, error)
	ImportOrders(path string) (int, error)
	ListReturns(page, limit int) ([]*models.Order, error)
	Process(userID, orderID, action string) error
	Return(orderID string) error
	Scroll(userID, lastID string, limit int) ([]*models.Order, string, error)
}

type orderService struct {
	storage storage.Storage
}

func New(storage storage.Storage) Service {
	return &orderService{storage: storage}
}
