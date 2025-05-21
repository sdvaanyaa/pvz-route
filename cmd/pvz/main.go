package main

import (
	"gitlab.ozon.dev/sd_vaanyaa/homework/cmd/commands"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/storage/json_storage"
	"log"
)

const storagePath = "data"

func main() {
	storage, err := json_storage.New(storagePath)
	if err != nil {
		log.Fatal(err)
	}

	orderService := services.New(storage)
	commands.Setup(orderService)
	commands.Execute()
}
