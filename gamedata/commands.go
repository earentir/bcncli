// File: gamedata/commands.go
package gamedata

import (
	"fmt"

	"bcncli/client"
	"bcncli/common"

	"github.com/spf13/cobra"
)

// Cmd is the root command for game data operations
var Cmd = &cobra.Command{
	Use:   "gamedata",
	Short: "Retrieve static game data",
}

func init() {
	// Add --cache flag to items command
	itemsCmd.Flags().BoolP("cache", "c", false, "Save output to gamedata-items.json")
	Cmd.AddCommand(itemsCmd)
}

// itemsCmd fetches global item data and optionally caches it
var itemsCmd = &cobra.Command{
	Use:   "items",
	Short: "Fetch item data",
	Run: func(cmd *cobra.Command, args []string) {
		payload := map[string]interface{}{"type": "itemData"}
		data := client.FetchDataOrExit(payload)

		// Check cache flag
		cache, _ := cmd.Flags().GetBool("cache")
		fileName := "itemid.json"
		common.LoadItemData(fileName, 3600, cache)

		fmt.Printf("Data cached to %s\n", fileName)
		common.PrintJSON(data)

	},
}
