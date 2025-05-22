package commands

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services"
	"log"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "pvz",
	Short: "PVZ CLI for managing orders",
}

func Setup(orderService *services.OrderService) {
	SetupAcceptCmd(orderService)
	SetupReturnCmd(orderService)
	SetupProcessCmd(orderService)
}

func Execute() {
	fmt.Println("PVZ CLI. Type 'help' for commands or 'exit' to quit.")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "exit" {
			break
		}

		if input == "" {
			continue
		}

		args := strings.Fields(input)
		os.Args = []string{os.Args[0]}
		os.Args = append(os.Args, args...)

		if err := rootCmd.Execute(); err != nil {
			log.Printf("cannot execute root command: %s", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("cannot read standard input: %s", err)
	}
}
