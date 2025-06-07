package httpgateway

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gitlab.ozon.dev/sd_vaanyaa/homework/api/gen"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/middleware/httpmw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

const defaultRPS = 5

// Gateway is an HTTP gateway for the gRPC service.
type Gateway struct {
	mux    *runtime.ServeMux
	server *http.Server
}

// New creates a new HTTP gateway.
func New(grpcAddr, httpAddr string) (*Gateway, error) {
	mux := runtime.NewServeMux()

	err := gen.RegisterOrderServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		grpcAddr,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	)

	if err != nil {
		return nil, err
	}

	handler := httpmw.RateLimiter(defaultRPS)(mux)

	return &Gateway{
		mux: mux,
		server: &http.Server{
			Addr:    httpAddr,
			Handler: handler,
		},
	}, nil
}

// Run runs the HTTP gateway server.
func (g *Gateway) Run() error {
	log.Printf("HTTP gateway server listening at %s", g.server.Addr)
	return g.server.ListenAndServe()
}

// Stop gracefully shuts down the HTTP gateway server.
func (g *Gateway) Stop(ctx context.Context) error {
	log.Printf("stopping HTTP gateway...")
	return g.server.Shutdown(ctx)
}
