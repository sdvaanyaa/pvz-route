package main

import (
	"gitlab.ozon.dev/sd_vaanyaa/homework/cmd/commands"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/storage/json_storage"
	"log"
)

const storagePath = "data"

func main() {
	storage, err := json_storage.New(storagePath)
	if err != nil {
		log.Fatal(err)
	}

	orderService := order.New(storage)
	commands.Setup(orderService)
	commands.Execute()
}
