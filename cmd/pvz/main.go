package main

import (
	"context"
	"errors"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/api/grpcserver"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/api/httpgateway"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/storage/jsonstorage"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

const (
	storagePath = "data"
	grpcAddress = "localhost:50051"
	httpAddress = "localhost:8080"
)

func main() {
	storage, err := jsonstorage.New(storagePath)
	if err != nil {
		log.Fatalf("failed to init json storage: %v", err)
	}
	orderService := order.New(storage)

	grpcServer := grpcserver.New(orderService)
	httpGateway, err := httpgateway.New(grpcAddress, httpAddress)
	if err != nil {
		log.Fatalf("failed to init http gateway: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err = grpcServer.Run(ctx, grpcAddress); err != nil {
			log.Printf("gRPC server failed: %v", err)
		}
	}()

	go func() {
		if err = httpGateway.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP gateway failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down servers...")

	if err = httpGateway.Stop(context.Background()); err != nil {
		log.Printf("HTTP gateway shutdown error: %v", err)
	}

	log.Println("servers stopped")

	//cli.Setup(orderService)
	//cli.Execute()
}
