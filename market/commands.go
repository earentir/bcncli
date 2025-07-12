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

// Item represents an entry from itemid.json, with every field included.
type Item struct {
	Name        string        `json:"name"`
	Emoji       string        `json:"emoji"`
	IDName      string        `json:"idName"`
	Uncraftable bool          `json:"uncraftable"`
	Attributes  []string      `json:"attributes"`
	LootSources []string      `json:"lootSources"`
	Recipe      []interface{} `json:"recipe"`
	ID          int           `json:"id"`
	FlatID      string        `json:"flatId"`
	Cost        int64         `json:"cost"`
	UsedToCraft []interface{} `json:"usedToCraft"`
	ImageURL    string        `json:"imageUrl"`
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
		debug, _ := cmd.Flags().GetBool("debug")
		sortField, _ := cmd.Flags().GetString("sort")

		// 1) fetch raw JSON
		payload := map[string]interface{}{"type": "marketPreview"}
		raw := client.FetchDataOrExit(payload)

		// 2) if --debug, just dump JSON
		if debug {
			client.PrintJSON(raw)
			return
		}

		// 3) unmarshal into our struct
		var resp OverviewResponse
		if err := json.Unmarshal(raw, &resp); err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse overview: %v\n", err)
			os.Exit(1)
		}

		// 4) load items and build name lookup
		items := loadItems()
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
		for key, val := range resp.Data {
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
		fmt.Fprintln(w, "ITEM\tVALUE")
		for _, r := range rows {
			fmt.Fprintf(w, "%s (%d)\t%s\n", r.Name, r.ID, formatPrice(r.Value))
		}
		w.Flush()
	},
}

// itemCmd lists market listings for a specific item
var itemCmd = &cobra.Command{
	Use:   "item [itemId]",
	Short: "List market listings for an item",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		itemId := client.ParseID(args[0])
		payload := map[string]interface{}{"type": "marketListings", "itemId": itemId}
		raw := client.FetchDataOrExit(payload)

		if debug {
			client.PrintJSON(raw)
			return
		}

		var listings []Listing
		if err := json.Unmarshal(raw, &listings); err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse listings: %v\n", err)
			os.Exit(1)
		}

		items := loadItems()
		nameByID := make(map[int]string, len(items))
		for _, it := range items {
			nameByID[it.ID] = it.Name
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "ITEM\tPRICE\tAMOUNT")
		for _, l := range listings {
			name := nameByID[l.ItemID]
			if name == "" {
				name = fmt.Sprintf("UNKNOWN(%d)", l.ItemID)
			}
			fmt.Fprintf(w, "%s (%d)\t%s\t%d\n", name, l.ItemID, formatPrice(l.Price), l.Amount)
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
		bcId := client.ParseID(args[0])
		payload := map[string]interface{}{"type": "userMarketListings", "id": bcId}
		raw := client.FetchDataOrExit(payload)

		// 3) If debug, just dump the raw JSON:
		if debug {
			client.PrintJSON(raw)
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
		var items []Item
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
			fmt.Fprintf(w, "%s (%d)\t%s\t%d\n", name, l.ItemID, formatPrice(l.Price), l.Amount)
		}
		w.Flush()
	},
}

// loadItems reads itemid.json next to the binary and returns all items.
func loadItems() []Item {
	exePath, _ := os.Executable()
	jsonPath := filepath.Join(filepath.Dir(exePath), "itemid.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read itemid.json: %v\n", err)
		os.Exit(1)
	}
	var items []Item
	if err := json.Unmarshal(data, &items); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse itemid.json: %v\n", err)
		os.Exit(1)
	}
	return items
}

// formatPrice formats n either with units (K, M, B, T) or, if you pass
// useNum=true, as a plain integer with space separators.
//
//	formatNumber(1230000)        == "1.23M"
//	formatNumber(1000000)        == "1M"
//	formatNumber(200000000)      == "200M"
//	formatNumber(123000000000)   == "123B"
//	formatNumber(123000, true)   == "123 000"
func formatPrice(n int64, useNum ...bool) string {
	num := false
	if len(useNum) > 0 && useNum[0] {
		num = true
	}

	// Plain number with spaces
	if num {
		s := strconv.FormatInt(n, 10)
		neg := strings.HasPrefix(s, "-")
		if neg {
			s = s[1:]
		}
		// group digits in threes
		var b strings.Builder
		pre := len(s) % 3
		if pre == 0 {
			pre = 3
		}
		b.WriteString(s[:pre])
		for i := pre; i < len(s); i += 3 {
			b.WriteByte(' ')
			b.WriteString(s[i : i+3])
		}
		out := b.String()
		if neg {
			out = "-" + out
		}
		return out
	}

	// Unit thresholds
	abs := n
	if abs < 0 {
		abs = -abs
	}
	units := []struct {
		thresh int64
		suf    string
	}{
		{1e12, "T"},
		{1e9, "B"},
		{1e6, "M"},
		{1e3, "K"},
	}
	for _, u := range units {
		if abs >= u.thresh {
			v := float64(n) / float64(u.thresh)
			// two decimals, then trim trailing zeros and dot
			str := strconv.FormatFloat(v, 'f', 2, 64)
			str = strings.TrimRight(str, "0")
			str = strings.TrimRight(str, ".")
			return str + u.suf
		}
	}

	// smaller than 1K: just print it
	return strconv.FormatInt(n, 10)
}
