package main

import (
	"context"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/api/grpcserver"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/storage/jsonstorage"
	"log"
)

const (
	storagePath = "data"
	grpcAddress = "localhost:50053" //TODO
)

func main() {
	storage, err := jsonstorage.New(storagePath)
	if err != nil {
		log.Fatalf("failed to init json storage: %v", err)
	}

	orderService := order.New(storage)

	grpcServer := grpcserver.New(orderService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err = grpcServer.Run(ctx, grpcAddress); err != nil {
		log.Fatalf("failed to run grpc server: %v", err)
	}

	//cli.Setup(orderService)
	//cli.Execute()
}
