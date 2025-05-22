package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services"
)

var listOrdersCmd = &cobra.Command{
	Use:   "list-orders",
	Short: "Get a list of orders",
}

func SetupListOrdersCmd(orderService *services.OrderService) {
	listOrdersCmd.Flags().StringP(FlagUserID, "u", "", "User ID")
	listOrdersCmd.Flags().BoolP(FlagInPVZ, "p", false, "Show only orders in PV")
	listOrdersCmd.Flags().IntP(FlagLast, "l", 0, "Show last N orders")
	listOrdersCmd.Flags().IntP(FlagPage, "n", 1, "Page number")
	listOrdersCmd.Flags().IntP(FlagLimit, "c", 0, "Orders per page")

	_ = listOrdersCmd.MarkFlagRequired(FlagUserID)

	listOrdersCmd.Run = func(cmd *cobra.Command, args []string) {
		userID, err := GetFlagString(cmd, FlagUserID)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}

		inPVZ, err := GetFlagBool(cmd, FlagInPVZ)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}

		last, err := GetFlagInt(cmd, FlagLast)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}

		page, err := GetFlagInt(cmd, FlagPage)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}

		limit, err := GetFlagInt(cmd, FlagLimit)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}

		orders, total, err := orderService.ListOrders(userID, inPVZ, last, page, limit)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}

		for _, o := range orders {
			fmt.Printf("ORDER: %s %s %s %s\n", o.ID, o.UserID, o.Status, o.StorageExpire)
		}

		fmt.Printf("TOTAL: %d\n", total)
	}
}

func init() {
	rootCmd.AddCommand(listOrdersCmd)
}
