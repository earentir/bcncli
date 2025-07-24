// File: gamedata/commands.go
package gamedata

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

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
	Cmd.AddCommand(itemCmd)
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

// itemCmd fetches details for a specific item by ID or name
var itemCmd = &cobra.Command{
	Use:   "item [id or name]",
	Short: "Fetch details for a specific item",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		arg := args[0]
		fileName := "itemid.json"
		items, err := common.LoadItemData(fileName, 3600, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading item data: %v\n", err)
			os.Exit(1)
		}

		// Find the requested item
		item, err := findItem(items, arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		// Print item details in table form
		printItemDetails(item, items)
	},
}

// findItem looks up an item by numeric ID or by name/idName
func findItem(items []common.Item, arg string) (common.Item, error) {
	// Try parsing as integer ID
	if id, err := strconv.Atoi(arg); err == nil {
		for _, it := range items {
			if it.ID == id {
				return it, nil
			}
		}
	}
	// Otherwise, match by name or idName (case-insensitive)
	lower := strings.ToLower(arg)
	for _, it := range items {
		if strings.ToLower(it.Name) == lower || strings.ToLower(it.IDName) == lower {
			return it, nil
		}
	}
	return common.Item{}, fmt.Errorf("item %q not found", arg)
}

// sanitizeEmoji returns a displayable emoji or alias for the terminal
func sanitizeEmoji(e string) string {
	// Discord-style <:name:id> custom emoji; fall back to alias
	if strings.HasPrefix(e, "<:") && strings.HasSuffix(e, ">") {
		parts := strings.Split(e, ":")
		if len(parts) >= 2 {
			return ":" + parts[1] + ":"
		}
	}
	return e
}

// printItemDetails outputs all fields of an item, and displays recipe components
func printItemDetails(item common.Item, allItems []common.Item) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Field\tValue\n")
	fmt.Fprintf(w, "ID\t%d\n", item.ID)
	fmt.Fprintf(w, "Emoji\t%s\n", sanitizeEmoji(item.Emoji))
	fmt.Fprintf(w, "Name\t %s\n", item.Name)
	fmt.Fprintf(w, "Description\t%s\n", item.Description)
	fmt.Fprintf(w, "Uncraftable\t%t\n", item.Uncraftable)
	fmt.Fprintf(w, "Cost\t%d\n", item.Cost)
	fmt.Fprintf(w, "Attributes\t%v\n", item.Attributes)
	fmt.Fprintf(w, "LootSources\t%v\n", item.LootSources)
	fmt.Fprintf(w, "UseLimit\t%d\n", item.UseLimit)
	// Replace item IDs with names
	var usedNames []string
	for _, uid := range item.UsedToCraft {
		usedNames = append(usedNames, common.LookUpItemName(uid, allItems))
	}
	fmt.Fprintf(w, "UsedToCraft	%s\n", strings.Join(usedNames, ", "))
	fmt.Fprintf(w, "ImageURL\t%s\n", item.ImageURL)

	// Display recipe if present
	if len(item.Recipe) > 0 {
		fmt.Fprintf(w, "\nRecipe Components:\t\n")
		fmt.Fprintf(w, "Name\tCount\n")
		// For each ItemRecipe, find matching Item and print
		for _, ir := range item.Recipe {
			for _, it := range allItems {
				if it.ID == ir.ID {
					fmt.Fprintf(w, "%s\t%d\n", it.Name, ir.Count)
					break
				}
			}
		}
	}
	w.Flush()
}
