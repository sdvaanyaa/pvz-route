package grpcserver

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/sd_vaanyaa/homework/api/gen"
	middleware "gitlab.ozon.dev/sd_vaanyaa/homework/internal/middleware/grpc"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/packaging"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
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

	log.Printf("grpc server listening at %s", addr)

	return s.grpcServer.Serve(listener)
}

func (s *Server) Accept(ctx context.Context, req *gen.AcceptOrderRequest) (*gen.OrderResponse, error) {
	log.Printf(
		"Accept request: order_id = %s, user_id = %s, expires_at = %s",
		req.OrderId,
		req.UserId,
		req.ExpiresAt,
	)

	packageType := packaging.PackageNone
	if req.PackageType != nil {
		packageType = packageTypeToString(*req.PackageType)
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
		log.Printf("failed to accept order: %v", err)
		return nil, mapError(err)
	}

	log.Printf("accepted order successfully: order_id = %s, status = %s", o.ID, o.Status)
	return &gen.OrderResponse{
		OrderId: o.ID,
		Status:  statusToProto(o.Status),
	}, nil
}

func (s *Server) Return(ctx context.Context, req *gen.OrderIdRequest) (*gen.OrderResponse, error) {
	log.Printf("Return request: order_id = %s", req.OrderId)

	if err := s.orderSvc.Return(req.OrderId); err != nil {
		log.Printf("failed to return order: order_id = %s, error: %v", req.OrderId, err)
		return nil, mapError(err)
	}

	log.Printf("Order returned successfully: order_id = %s, status = archived", req.OrderId)

	return &gen.OrderResponse{
		OrderId: req.OrderId,
		Status:  gen.OrderStatus_ORDER_STATUS_ARCHIVED,
	}, nil
}
