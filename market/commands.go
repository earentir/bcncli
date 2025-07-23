// Package market provides commands to interact with the marketplace.
package market

import (
	"bcncli/client"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"bcncli/common"

	"github.com/spf13/cobra"
)

// Listing represents a single market listing from the API.
type Listing struct {
	ID     int64 `json:"id"`
	BcID   int64 `json:"bcId"`
	ItemID int   `json:"itemId"`
	Price  int64 `json:"price"`
	Amount int64 `json:"amount"`
}

// OverviewResponse models the marketPreview API response.
type OverviewResponse struct {
	LastUpdated int64            `json:"lastUpdated"`
	Data        map[string]int64 `json:"data"`
}

// Cmd is the root command for market operations
var Cmd = &cobra.Command{
	Use:   "market",
	Short: "Manage marketplace operations",
}

func init() {
	overviewCmd.Flags().BoolP("debug", "d", false, "print raw JSON response")
	overviewCmd.Flags().StringP("sort", "s", "id", "sort overview by: id, name, or price")
	itemCmd.Flags().BoolP("debug", "d", false, "print raw JSON response")
	userCmd.Flags().BoolP("debug", "d", false, "print raw JSON response")
	Cmd.AddCommand(overviewCmd, itemCmd, userCmd)
}

var overviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "Show market overview",
	Run: func(cmd *cobra.Command, args []string) {
		sortField, _ := cmd.Flags().GetString("sort")

		// 1) fetch raw JSON
		payload := map[string]interface{}{"type": "marketPreview"}
		raw := client.FetchDataOrExit(payload)

		// 2) if --debug, just dump JSON
		if debug, _ := cmd.Flags().GetBool("debug"); debug {
			common.PrintJSON(raw)
			return
		}

		// 3) unmarshal into our struct
		var responce OverviewResponse
		if err := json.Unmarshal(raw, &responce); err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse overview: %v\n", err)
			os.Exit(1)
		}

		// 4) load items and build name lookup
		items, err := common.LoadItemData("itemid.json", 3600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not load items: %v\n", err)
			os.Exit(1)
		}

		nameByID := make(map[int]string, len(items))
		for _, it := range items {
			nameByID[it.ID] = it.Name
		}

		// 5) build rows slice
		type row struct {
			ID    int
			Name  string
			Value int64
		}
		var rows []row
		for key, val := range responce.Data {
			if !strings.HasPrefix(key, "item") {
				continue
			}
			idNum, err := strconv.Atoi(strings.TrimPrefix(key, "item"))
			if err != nil {
				continue
			}
			name := nameByID[idNum]
			if name == "" {
				name = fmt.Sprintf("UNKNOWN(%d)", idNum)
			}
			rows = append(rows, row{ID: idNum, Name: name, Value: val})
		}

		// 6) sort according to --sort
		switch strings.ToLower(sortField) {
		case "id":
			sort.Slice(rows, func(i, j int) bool { return rows[i].ID < rows[j].ID })
		case "name":
			sort.Slice(rows, func(i, j int) bool { return rows[i].Name < rows[j].Name })
		case "price":
			sort.Slice(rows, func(i, j int) bool { return rows[i].Value < rows[j].Value })
		default:
			fmt.Fprintf(os.Stderr, "invalid sort option: %s (must be id, name, or price)\n", sortField)
			os.Exit(1)
		}

		// 7) render as table
		w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
		defer w.Flush()

		fmt.Fprintln(w, "ITEM\tVALUE")
		for _, r := range rows {
			fmt.Fprintf(w, "%s (%d)\t%s\n", r.Name, r.ID, common.FormatPrice(r.Value))
		}
		w.Flush()
	},
}

var itemCmd = &cobra.Command{
	Use:   "item [itemId]",
	Short: "List market listings for an item",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		itemID := common.ParseID(args[0])

		// instead of a map, build a typed payload
		type requestPayload struct {
			Type   string `json:"type"`
			ItemID int    `json:"itemId"`
		}
		payload := map[string]interface{}{"type": "marketListings", "itemId": itemID}
		raw := client.FetchDataOrExit(payload)

		if debug {
			common.PrintJSON(raw)
			return
		}

		// unmarshal into a slice of Listing structs
		var listings []Listing
		if err := json.Unmarshal(raw, &listings); err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse listings: %v\n", err)
			os.Exit(1)
		}

		// load item names as before
		items, err := common.LoadItemData("itemid.json", 3600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not load items: %v\n", err)
			os.Exit(1)
		}

		nameByID := make(map[int]string, len(items))
		for _, it := range items {
			nameByID[it.ID] = it.Name
		}

		// pretty-print
		w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "ITEM\tPRICE\tAMOUNT")
		for _, l := range listings {
			name := nameByID[l.ItemID]
			if name == "" {
				name = fmt.Sprintf("UNKNOWN(%d)", l.ItemID)
			}
			fmt.Fprintf(
				w,
				"%s (%d)\t%s\t%d\n",
				name,                        // item name
				l.ItemID,                    // item ID
				common.FormatPrice(l.Price), // price in your format
				l.Amount,                    // quantity
			)
		}
		w.Flush()
	},
}

// userCmd lists a user's market listings
var userCmd = &cobra.Command{
	Use:   "user [bcId]",
	Short: "List market listings for a user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 1) Parse flag:
		debug, _ := cmd.Flags().GetBool("debug")

		// 2) Fetch raw data:
		bcID := common.ParseID(args[0])
		payload := map[string]interface{}{"type": "userMarketListings", "id": bcID}
		raw := client.FetchDataOrExit(payload)

		// 3) If debug, just dump the raw JSON:
		if debug {
			common.PrintJSON(raw)
			return
		}

		// 4) Otherwise unmarshal into []Listing
		var listings []Listing
		if err := json.Unmarshal(raw, &listings); err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse listings: %v\n", err)
			os.Exit(1)
		}

		// 5) Load item definitions from itemid.json:
		exePath, _ := os.Executable()
		jsonPath := filepath.Join(filepath.Dir(exePath), "itemid.json")
		b, err := os.ReadFile(jsonPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read itemid.json: %v\n", err)
			os.Exit(1)
		}
		var items []common.Item
		if err := json.Unmarshal(b, &items); err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse itemid.json: %v\n", err)
			os.Exit(1)
		}
		// Build map[id]name
		nameByID := make(map[int]string, len(items))
		for _, it := range items {
			nameByID[it.ID] = it.Name
		}

		// 6) Print a nice table:
		w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "ITEM\tPRICE\tAMOUNT")
		for _, l := range listings {
			name, ok := nameByID[l.ItemID]
			if !ok {
				name = fmt.Sprintf("UNKNOWN(%d)", l.ItemID)
			}
			fmt.Fprintf(w, "%s (%d)\t%s\t%d\n", name, l.ItemID, common.FormatPrice(l.Price), l.Amount)
		}
		w.Flush()
	},
}
