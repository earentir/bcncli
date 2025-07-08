// File: gamedata/commands.go
package gamedata

import (
	"fmt"
	"os"

	"bcncli/client"

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
		if cache {
			fileName := "gamedata-items.json"
			if err := os.WriteFile(fileName, data, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing cache file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Data cached to %s\n", fileName)
		} else {
			client.PrintJSON(data)
		}
	},
}
