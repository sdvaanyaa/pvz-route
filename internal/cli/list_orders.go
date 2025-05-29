package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/packaging"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
	"time"
)

var listOrdersCmd = &cobra.Command{
	Use:   "list-orders",
	Short: "Get a list of orders",
}

func setupListOrdersCmd(orderSvc order.Service) {
	listOrdersCmd.Flags().StringP(FlagUserID, ShortUserID, "", "User ID")
	listOrdersCmd.Flags().BoolP(FlagInPVZ, ShortInPVZ, false, "Show only orders in PV")
	listOrdersCmd.Flags().IntP(FlagLast, ShortLast, 0, "Show last N orders")
	listOrdersCmd.Flags().IntP(FlagPage, ShortPage, 1, "Page number")
	listOrdersCmd.Flags().IntP(FlagLimit, ShortLimit, 0, "Orders per page")

	_ = listOrdersCmd.MarkFlagRequired(FlagUserID)

	listOrdersCmd.Run = func(cmd *cobra.Command, args []string) {
		userID, inPVZ, last, page, limit, err := parseListOrdersFlags(cmd)
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
			formatedTime := o.StorageExpire.Format(time.DateOnly)

			if o.PackageType == "" {
				o.PackageType = packaging.PackageNone
			}

			fmt.Printf(
				"ORDER: %s %s %s %s %s %.2f %.2f\n",
				o.ID,
				o.UserID,
				o.Status,
				formatedTime,
				o.PackageType,
				o.Weight,
				o.Price,
			)
		}

		fmt.Printf("TOTAL: %d\n", total)
	}
}

func init() {
	rootCmd.AddCommand(listOrdersCmd)
}

func parseListOrdersFlags(cmd *cobra.Command) (string, bool, int, int, int, error) {
	userID, err := getFlagString(cmd, FlagUserID)
	if err != nil {
		return "", false, 0, 0, 0, err
	}

	inPVZ, err := getFlagBool(cmd, FlagInPVZ)
	if err != nil {
		return "", false, 0, 0, 0, err
	}

	last, err := getFlagInt(cmd, FlagLast)
	if err != nil {
		return "", false, 0, 0, 0, err
	}

	page, err := getFlagInt(cmd, FlagPage)
	if err != nil {
		return "", false, 0, 0, 0, err
	}

	limit, err := getFlagInt(cmd, FlagLimit)
	if err != nil {
		return "", false, 0, 0, 0, err
	}

	return userID, inPVZ, last, page, limit, nil
}
