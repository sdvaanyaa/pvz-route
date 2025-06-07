package grpcserver

import (
	"errors"
	"gitlab.ozon.dev/sd_vaanyaa/homework/api/gen"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/models"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/packaging"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func protoPackageTypeToString(pt gen.PackageType) string {
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

func stringPackageTypeToProto(pt string) gen.PackageType {
	switch pt {
	case packaging.PackageNone:
		return gen.PackageType_PACKAGE_TYPE_NONE
	case packaging.PackageBag:
		return gen.PackageType_PACKAGE_TYPE_BAG
	case packaging.PackageBox:
		return gen.PackageType_PACKAGE_TYPE_BOX
	case packaging.PackageFilm:
		return gen.PackageType_PACKAGE_TYPE_FILM
	case packaging.PackageBagFilm:
		return gen.PackageType_PACKAGE_TYPE_BAG_FILM
	case packaging.PackageBoxFilm:
		return gen.PackageType_PACKAGE_TYPE_BOX_FILM
	default:
		return gen.PackageType_PACKAGE_TYPE_NONE
	}
}

func protoActionToString(at gen.ActionType) string {
	switch at {
	case gen.ActionType_ACTION_TYPE_ISSUE:
		return "issue"
	case gen.ActionType_ACTION_TYPE_RETURN:
		return "return"
	default:
		return "invalid action type"
	}
}

func stringStatusToProto(status string) gen.OrderStatus {
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

func modelsOrderToProto(o *models.Order) *gen.Order {
	packageType := gen.PackageType_PACKAGE_TYPE_NONE
	if o.PackageType != "" {
		packageType = stringPackageTypeToProto(o.PackageType)
	}

	var issuedAt, returnedAt, archivedAt *timestamppb.Timestamp
	if o.IssuedAt != nil && !o.IssuedAt.IsZero() {
		issuedAt = timestamppb.New(*o.IssuedAt)
	}

	if o.ReturnedAt != nil && !o.ReturnedAt.IsZero() {
		returnedAt = timestamppb.New(*o.ReturnedAt)
	}

	if o.ArchivedAt != nil && !o.ArchivedAt.IsZero() {
		archivedAt = timestamppb.New(*o.ArchivedAt)
	}

	return &gen.Order{
		OrderId:     o.ID,
		UserId:      o.UserID,
		Status:      stringStatusToProto(o.Status),
		ExpiresAt:   timestamppb.New(o.StorageExpire),
		Weight:      o.Weight,
		TotalPrice:  o.Price,
		PackageType: &packageType,
		CreatedAt:   timestamppb.New(o.CreatedAt),
		IssuedAt:    issuedAt,
		ReturnedAt:  returnedAt,
		ArchivedAt:  archivedAt,
	}
}

func toGRPCError(err error) error {
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
