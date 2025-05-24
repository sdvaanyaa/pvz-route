package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
)

var listOrdersCmd = &cobra.Command{
	Use:   "list-orders",
	Short: "Get a list of orders",
}

func setupListOrdersCmd(orderSvc order.Service) {
	listOrdersCmd.Flags().StringP(FlagUserID, "u", "", "User ID")
	listOrdersCmd.Flags().BoolP(FlagInPVZ, "p", false, "Show only orders in PV")
	listOrdersCmd.Flags().IntP(FlagLast, "l", 0, "Show last N orders")
	listOrdersCmd.Flags().IntP(FlagPage, "n", 1, "Page number")
	listOrdersCmd.Flags().IntP(FlagLimit, "m", 0, "Orders per page")

	_ = listOrdersCmd.MarkFlagRequired(FlagUserID)

	listOrdersCmd.Run = func(cmd *cobra.Command, args []string) {
		userID, err := getFlagString(cmd, FlagUserID)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		inPVZ, err := getFlagBool(cmd, FlagInPVZ)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		last, err := getFlagInt(cmd, FlagLast)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		page, err := getFlagInt(cmd, FlagPage)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		limit, err := getFlagInt(cmd, FlagLimit)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		orders, total, err := orderSvc.ListOrders(userID, inPVZ, last, page, limit)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		for _, o := range orders {
			formatedTime := o.StorageExpire.Format("02 Jan 15:04")
			fmt.Printf("ORDER: %s %s %s %s\n", o.ID, o.UserID, o.Status, formatedTime)
		}

		fmt.Printf("TOTAL: %d\n", total)
	}
}

func init() {
	rootCmd.AddCommand(listOrdersCmd)
}
