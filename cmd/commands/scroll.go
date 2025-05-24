package commands

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
	"os"
	"strings"
)

var scrollCmd = &cobra.Command{
	Use:   "scroll-orders",
	Short: "Fetch user orders with infinite scrolling support",
}

func SetupScrollCmd(orderSvc order.Service) {
	scrollCmd.Flags().StringP(FlagUserID, "u", "", "User ID")
	scrollCmd.Flags().StringP(FlagLast, "l", "0", "Last ID")
	scrollCmd.Flags().IntP(FlagLimit, "n", 20, "Number of orders to fetch (default 20)")

	_ = scrollCmd.MarkFlagRequired(FlagUserID)

	scrollCmd.Run = func(cmd *cobra.Command, args []string) {
		userID, err := GetFlagString(cmd, FlagUserID)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		lastID, err := GetFlagString(cmd, FlagLast)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		limit, err := GetFlagInt(cmd, FlagLimit)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		if err = runScrollLoop(orderSvc, userID, lastID, limit); err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}
	}
}

func runScrollLoop(svc order.Service, userID, lastID string, limit int) error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		orders, nextLastID, err := svc.Scroll(userID, lastID, limit)
		if err != nil {
			return err
		}

		for _, o := range orders {
			fmt.Printf("ORDER: %s %s %s\n",
				o.ID,
				o.Status,
				o.StorageExpire.Format("02 Jan 15:04"),
			)
		}

		fmt.Printf("NEXT: %s\n", nextLastID)

		if nextLastID == "" {
			fmt.Println("No more orders. Exiting scroll mode.")
			return nil
		}

		fmt.Print("> ")
		if !scanner.Scan() {
			return fmt.Errorf("failed to read input: %w", scanner.Err())
		}
		input := strings.TrimSpace(scanner.Text())

		switch input {
		case "next":
			lastID = nextLastID
		case "exit":
			return nil
		default:
			fmt.Println("expected 'next' or 'exit'")
		}
	}
}

func init() {
	rootCmd.AddCommand(scrollCmd)
}
