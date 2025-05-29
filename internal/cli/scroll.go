package cli

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/packaging"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
	"os"
	"strings"
	"time"
)

var scrollCmd = &cobra.Command{
	Use:   "scroll-orders",
	Short: "Fetch user orders with infinite scrolling support",
}

func setupScrollCmd(orderSvc order.Service) {
	scrollCmd.Flags().StringP(FlagUserID, ShortUserID, "", "User ID")
	scrollCmd.Flags().StringP(FlagLast, ShortLast, "0", "Last ID")
	scrollCmd.Flags().IntP(FlagLimit, ShortLimit, 20, "Number of orders to fetch (default 20)")

	_ = scrollCmd.MarkFlagRequired(FlagUserID)

	scrollCmd.Run = func(cmd *cobra.Command, args []string) {
		userID, lastID, limit, err := parseScrollFlags(cmd)
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

func parseScrollFlags(cmd *cobra.Command) (string, string, int, error) {
	userID, err := getFlagString(cmd, FlagUserID)
	if err != nil {
		return "", "", 0, err
	}

	lastID, err := getFlagString(cmd, FlagLast)
	if err != nil {
		return "", "", 0, err
	}

	limit, err := getFlagInt(cmd, FlagLimit)
	if err != nil {
		return "", "", 0, err
	}

	return userID, lastID, limit, nil
}

func runScrollLoop(svc order.Service, userID, lastID string, limit int) error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		orders, nextLastID, err := svc.Scroll(userID, lastID, limit)
		if err != nil {
			return err
		}

		for _, o := range orders {
			formatedTime := o.StorageExpire.Format(time.DateOnly)

			if o.PackageType == "" {
				o.PackageType = packaging.PackageNone
			}

			fmt.Printf("ORDER: %s %s %s %s %.2f %.2f\n",
				o.ID,
				o.Status,
				formatedTime,
				o.PackageType,
				o.Weight,
				o.Price,
			)
		}

		fmt.Printf("NEXT: %s\n", nextLastID)

		if nextLastID == "" {
			fmt.Println("no more orders, exiting scroll mode.")
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
