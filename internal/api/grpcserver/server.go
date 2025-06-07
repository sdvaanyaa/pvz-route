package grpcserver

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/sd_vaanyaa/homework/api/gen"
	middleware "gitlab.ozon.dev/sd_vaanyaa/homework/internal/middleware/grpcmw"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/packaging"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
)

const (
	defaultLastN = 0
	defaultPage  = 1
	defaultLimit = 0
	defaultInPVZ = false
)

type Server struct {
	gen.UnimplementedOrderServiceServer
	orderSvc   order.Service
	grpcServer *grpc.Server
}

// New creates a new gRPC server.
func New(orderSvc order.Service) *Server {
	server := &Server{
		orderSvc: orderSvc,
	}

	server.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.LoggingInterceptor(),
			middleware.ValidationInterceptor(),
		),
	)
	gen.RegisterOrderServiceServer(server.grpcServer, server)

	reflection.Register(server.grpcServer)

	return server
}

// Run starts the gRPC server on the given address.
func (s *Server) Run(ctx context.Context, addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	log.Printf("grpcmw server listening at %s", addr)

	return s.grpcServer.Serve(listener)
}

func (s *Server) Accept(ctx context.Context, req *gen.AcceptOrderRequest) (*gen.OrderResponse, error) {
	packageType := packaging.PackageNone
	if req.PackageType != nil {
		packageType = protoPackageTypeToString(*req.PackageType)
	}

	o, err := s.orderSvc.Accept(
		req.OrderId,
		req.UserId,
		req.ExpiresAt.AsTime().Format("2006-01-02"),
		req.Weight,
		req.Price,
		packageType,
	)

	if err != nil {
		return nil, toGRPCError(err)
	}

	return &gen.OrderResponse{
		OrderId: o.ID,
		Status:  stringStatusToProto(o.Status),
	}, nil
}

func (s *Server) Return(ctx context.Context, req *gen.OrderIdRequest) (*gen.OrderResponse, error) {
	if err := s.orderSvc.Return(req.OrderId); err != nil {
		return nil, toGRPCError(err)
	}

	return &gen.OrderResponse{
		OrderId: req.OrderId,
		Status:  gen.OrderStatus_ORDER_STATUS_ARCHIVED,
	}, nil
}

func (s *Server) Process(ctx context.Context, req *gen.ProcessOrdersRequest) (*gen.ProcessResult, error) {
	result := &gen.ProcessResult{}

	for _, orderID := range req.OrderIds {
		if err := s.orderSvc.Process(req.UserId, orderID, protoActionToString(req.Action)); err != nil {
			result.Errors = append(result.Errors, orderID)
			continue
		}
		result.Processed = append(result.Processed, orderID)
	}

	return result, nil
}

func (s *Server) ListOrders(ctx context.Context, req *gen.ListOrdersRequest) (*gen.OrdersList, error) {
	lastN := defaultLastN
	if req.LastN != nil {
		lastN = int(*req.LastN)
	}

	page, limit := defaultPage, defaultLimit
	if req.Pagination != nil {
		page = int(req.Pagination.Page)
		limit = int(req.Pagination.CountOnPage)
	}

	inPVZ := defaultInPVZ
	if req.InPvz != nil {
		inPVZ = *req.InPvz
	}

	orders, total, err := s.orderSvc.ListOrders(req.UserId, inPVZ, lastN, page, limit)
	if err != nil {
		return nil, toGRPCError(err)
	}

	resp := &gen.OrdersList{
		Total: int32(total),
	}
	for _, o := range orders {
		resp.Orders = append(resp.Orders, modelsOrderToProto(o))
	}

	return resp, nil
}

func (s *Server) ListReturns(ctx context.Context, req *gen.ListReturnsRequest) (*gen.ReturnsList, error) {
	page, limit := defaultPage, defaultLimit
	if req.Pagination != nil {
		page = int(req.Pagination.Page)
		limit = int(req.Pagination.CountOnPage)
	}

	orders, err := s.orderSvc.ListReturns(page, limit)
	if err != nil {
		return nil, toGRPCError(err)
	}

	resp := &gen.ReturnsList{}
	for _, o := range orders {
		resp.Returns = append(resp.Returns, modelsOrderToProto(o))
	}

	return resp, nil
}

func (s *Server) History(ctx context.Context, _ *gen.GetHistoryRequest) (*gen.OrderHistoryList, error) {
	entries, err := s.orderSvc.History()
	if err != nil {
		return nil, toGRPCError(err)
	}

	resp := &gen.OrderHistoryList{}
	for _, entry := range entries {
		resp.History = append(resp.History, &gen.OrderHistory{
			OrderId:   entry.OrderID,
			Status:    stringStatusToProto(entry.Status),
			CreatedAt: timestamppb.New(entry.Timestamp),
		})
	}

	return resp, nil
}
