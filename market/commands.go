package market

import (
	"bcncli/client"

	"github.com/spf13/cobra"
)

// Cmd is the root command for market operations
var Cmd = &cobra.Command{
	Use:   "market",
	Short: "Manage marketplace operations",
}

func init() {
	Cmd.AddCommand(overviewCmd, itemCmd, userCmd)
}

// overviewCmd fetches a market preview
var overviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "Show market overview",
	Run: func(cmd *cobra.Command, args []string) {
		payload := map[string]interface{}{"type": "marketPreview"}
		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}

// itemCmd lists market listings for a specific item
var itemCmd = &cobra.Command{
	Use:   "item [itemId]",
	Short: "List market listings for an item",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		itemId := client.ParseID(args[0])
		payload := map[string]interface{}{"type": "marketListings", "itemId": itemId}
		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}

// userCmd lists a user's market listings
var userCmd = &cobra.Command{
	Use:   "user [bcId]",
	Short: "List market listings for a user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bcId := client.ParseID(args[0])
		payload := map[string]interface{}{"type": "userMarketListings", "id": bcId}
		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}
