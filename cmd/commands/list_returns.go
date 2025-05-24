package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
)

var listReturnsCmd = &cobra.Command{
	Use:   "list-returns",
	Short: "Get a list of returns",
}

func SetupListReturnsCmd(orderSvc order.Service) {
	listReturnsCmd.Flags().IntP(FlagPage, "n", 1, "Page number")
	listReturnsCmd.Flags().IntP(FlagLimit, "m", 0, "Orders per page")

	_ = listReturnsCmd.MarkFlagRequired(FlagPage)
	_ = listReturnsCmd.MarkFlagRequired(FlagLimit)

	listReturnsCmd.Run = func(cmd *cobra.Command, args []string) {
		page, err := GetFlagInt(cmd, FlagPage)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		limit, err := GetFlagInt(cmd, FlagLimit)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		orders, err := orderSvc.ListReturns(page, limit)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		for _, o := range orders {
			fmt.Printf("RETURN: %s %s %s\n", o.ID, o.UserID, o.ReturnedAt)
		}

		fmt.Printf("PAGE: %d LIMIT: %d\n", page, limit)
	}
}

func init() {
	rootCmd.AddCommand(listReturnsCmd)
}
