package grpcserver

import (
	"errors"
	"gitlab.ozon.dev/sd_vaanyaa/homework/api/gen"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/packaging"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func packageTypeToString(pt gen.PackageType) string {
	switch pt {
	case gen.PackageType_PACKAGE_TYPE_NONE:
		return packaging.PackageNone
	case gen.PackageType_PACKAGE_TYPE_BAG:
		return packaging.PackageBag
	case gen.PackageType_PACKAGE_TYPE_BOX:
		return packaging.PackageBox
	case gen.PackageType_PACKAGE_TYPE_FILM:
		return packaging.PackageFilm
	case gen.PackageType_PACKAGE_TYPE_BAG_FILM:
		return packaging.PackageBagFilm
	case gen.PackageType_PACKAGE_TYPE_BOX_FILM:
		return packaging.PackageBoxFilm
	default:
		return "invalid package type"
	}
}

func statusToProto(status string) gen.OrderStatus {
	switch status {
	case models.StatusAccepted:
		return gen.OrderStatus_ORDER_STATUS_ACCEPTED
	case models.StatusIssued:
		return gen.OrderStatus_ORDER_STATUS_ISSUED
	case models.StatusReturned:
		return gen.OrderStatus_ORDER_STATUS_RETURNED
	case models.StatusArchived:
		return gen.OrderStatus_ORDER_STATUS_ARCHIVED
	default:
		return gen.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

func mapError(err error) error {
	switch {
	case errors.Is(err, order.ErrEmptyOrderID),
		errors.Is(err, order.ErrEmptyUserID),
		errors.Is(err, order.ErrInvalidPageNumber),
		errors.Is(err, order.ErrInvalidLimitNumber),
		errors.Is(err, order.ErrInvalidDeadlineFormat),
		errors.Is(err, order.ErrUnknownAction),
		errors.Is(err, order.ErrInvalidPrice),
		errors.Is(err, order.ErrInvalidWeight),
		errors.Is(err, order.ErrEmptyFilePath),
		errors.Is(err, order.ErrEmptyImportFile),
		errors.Is(err, order.ErrEmptyValidOrders):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, order.ErrStorageNotExpired),
		errors.Is(err, order.ErrReturnPeriodExpired),
		errors.Is(err, order.ErrStorageExpired),
		errors.Is(err, order.ErrOrderIssued),
		errors.Is(err, order.ErrOrderNotIssued),
		errors.Is(err, order.ErrOrderNotAccepted):
		return status.Error(codes.FailedPrecondition, err.Error())

	case
		errors.Is(err, order.ErrOrderNotBelongsToUser):
		return status.Error(codes.PermissionDenied, err.Error())

	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
