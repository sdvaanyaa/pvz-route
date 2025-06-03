package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
	"strings"
)

var processCmd = &cobra.Command{
	Use:   "process-orders",
	Short: "Issue orders or accept customer returns",
}

func setupProcessCmd(orderSvc order.Service) {
	processCmd.Flags().StringP(FlagUserID, ShortUserID, "", "User ID")
	processCmd.Flags().StringP(FlagAction, ShortAction, "", "Action: issue or return")
	processCmd.Flags().StringP(FlagOrderIDs, ShortOrderID, "", "Order IDs")

	_ = processCmd.MarkFlagRequired(FlagUserID)
	_ = processCmd.MarkFlagRequired(FlagAction)
	_ = processCmd.MarkFlagRequired(FlagOrderIDs)

	processCmd.Run = func(cmd *cobra.Command, args []string) {
		userID, action, orderIDs, err := parseProcessFlags(cmd)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		sliceOrderIDs := strings.Split(orderIDs, ",")

		for _, orderID := range sliceOrderIDs {
			if err = orderSvc.Process(userID, orderID, action); err != nil {
				fmt.Printf("ERROR %s: %s\n", orderID, err)
				continue
			}

			fmt.Printf("PROCESSED: %s\n", orderID)
		}
	}
}

func init() {
	rootCmd.AddCommand(processCmd)
}

func parseProcessFlags(cmd *cobra.Command) (string, string, string, error) {
	userID, err := getFlagString(cmd, FlagUserID)
	if err != nil {
		return "", "", "", err
	}

	action, err := getFlagString(cmd, FlagAction)
	if err != nil {
		return "", "", "", err
	}

	orderIDs, err := getFlagString(cmd, FlagOrderIDs)
	if err != nil {
		return "", "", "", err
	}

	return userID, action, orderIDs, nil
}
