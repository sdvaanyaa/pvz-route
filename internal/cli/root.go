package cli

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
	"log"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Short: "PVZ CLI for managing orders",

	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

// Setup initializes all CLI commands with the provided order service.
func Setup(orderSvc order.Service) {
	setupAcceptCmd(orderSvc)
	setupReturnCmd(orderSvc)
	setupProcessCmd(orderSvc)
	setupListOrdersCmd(orderSvc)
	setupListReturnsCmd(orderSvc)
	setupHistoryCmd(orderSvc)
	setupImportCmd(orderSvc)
	setupScrollCmd(orderSvc)
}

// Execute starts an interactive CLI session, reading user input line by line.
// It prompts the user, parses the input into commands and arguments,
// resets flags from the previous command invocation before each new execution,
// and runs the corresponding command.
// The session continues until the user types "exit" or input ends.
// Any input reading errors are logged.
func Execute() {
	fmt.Println("PVZ CLI. Type 'help' for cli or 'exit' to quit.")

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

		for _, subCmd := range rootCmd.Commands() {
			resetFlags(subCmd)
		}

		args := strings.Fields(input)
		os.Args = []string{os.Args[0]}
		os.Args = append(os.Args, args...)

		_ = rootCmd.Execute()
	}

	if err := scanner.Err(); err != nil {
		log.Printf("cannot read standard input: %s", err)
	}
}
