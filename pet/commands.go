package pet

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"bcncli/client"
	"bcncli/common"

	"github.com/spf13/cobra"
)

// Pet represents the subset of fields to display in the table
type Pet struct {
	ID                 int64          `json:"id"`
	OwnerBCID          int64          `json:"ownerBcId"`
	HatchDate          string         `json:"hatchDate"`
	Name               string         `json:"name"`
	Tier               int            `json:"tier"`
	XP                 int64          `json:"xp"`
	Species            string         `json:"species"`
	Generation         int            `json:"generation"`
	ParentAID          int64          `json:"parentAId"`
	ParentBID          int64          `json:"parentBId"`
	TimesBred          int            `json:"timesBred"`
	LastBred           string         `json:"lastBred"`
	HeldItemID         int64          `json:"heldItemId"`
	UnsyncedEnergy     int64          `json:"unsyncedEnergy"`
	AdventureType      string         `json:"adventureType"`
	AdventureBoost     AdventureBoost `json:"adventureBoost"`
	LastAdventureSync  string         `json:"lastAdventureSync"`
	LifetimeItemsFound int64          `json:"lifetimeItemsFound"`
	Craving            Craving        `json:"craving"`
	Skin               string         `json:"skin"`
	Aura               string         `json:"aura"`
}

// AdventureBoost represents the boost details for a pet's adventure
type AdventureBoost struct {
	Multiplier int   `json:"multiplier"`
	EndTime    int64 `json:"endTime"`
}

// Craving represents the craving details for a pet
type Craving struct {
	ItemID int64 `json:"itemId"`
	Amount int64 `json:"amount"`
}

var Cmd = &cobra.Command{
	Use:   "pet",
	Short: "Manage pets",
}

func init() {
	// Register subcommands
	Cmd.AddCommand(infoCmd, ownedCmd, offspringCmd)

	// Define flags for owned command
	ownedCmd.Flags().Bool("debug", false, "Enable debug (JSON) output")
	ownedCmd.Flags().String("sort", "", "Sort table by this field (id, name, species, tier, xp, adventure, items, boost)")
	ownedCmd.Flags().String("group", "", "Group table by this field (species, tier, boost, adventure, craving)")
}

var infoCmd = &cobra.Command{
	Use:   "info [id]",
	Short: "Fetch pet info",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := common.ParseID(args[0])
		payload := map[string]interface{}{"type": "pet", "id": id}
		data := client.FetchDataOrExit(payload)
		common.PrintJSON(data)
	},
}

var ownedCmd = &cobra.Command{
	Use:   "owned [userId]",
	Short: "List pets for a user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userID := common.ParseID(args[0])
		payload := map[string]interface{}{"type": "userPetsAndEggs", "id": userID}
		raw := client.FetchDataOrExit(payload)

		// Parse wrapper
		var resp map[string]json.RawMessage
		if err := json.Unmarshal(raw, &resp); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing response wrapper: %v\n", err)
			os.Exit(1)
		}

		// Debug JSON
		if debug, _ := cmd.Flags().GetBool("debug"); debug {
			common.PrintJSON(resp["pets"])
			return
		}

		// Unmarshal pets
		var pets []Pet
		if err := json.Unmarshal(resp["pets"], &pets); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing pets data: %v\n", err)
			os.Exit(1)
		}

		// Sort
		sortKey, _ := cmd.Flags().GetString("sort")
		if sortKey != "" {
			switch sortKeyLower := strings.ToLower(sortKey); sortKeyLower {
			case "id":
				sort.Slice(pets, func(i, j int) bool { return pets[i].ID < pets[j].ID })
			case "name":
				sort.Slice(pets, func(i, j int) bool { return pets[i].Name < pets[j].Name })
			case "species":
				sort.Slice(pets, func(i, j int) bool { return pets[i].Species < pets[j].Species })
			case "tier":
				sort.Slice(pets, func(i, j int) bool { return pets[i].Tier < pets[j].Tier })
			case "xp":
				sort.Slice(pets, func(i, j int) bool { return pets[i].XP < pets[j].XP })
			case "adventure":
				sort.Slice(pets, func(i, j int) bool { return pets[i].AdventureType < pets[j].AdventureType })
			case "items":
				sort.Slice(pets, func(i, j int) bool { return pets[i].LifetimeItemsFound < pets[j].LifetimeItemsFound })
			case "boost":
				sort.Slice(pets, func(i, j int) bool { return pets[i].AdventureBoost.Multiplier < pets[j].AdventureBoost.Multiplier })
			default:
				fmt.Fprintf(os.Stderr, "Invalid sort field: %s\n", sortKey)
				os.Exit(1)
			}
		}

		// Group
		groupKey, _ := cmd.Flags().GetString("group")
		if groupKey != "" {
			// Build groups
			groups := make(map[string][]Pet)
			for _, p := range pets {
				var key string
				switch strings.ToLower(groupKey) {
				case "species":
					key = p.Species
				case "tier":
					key = fmt.Sprint(p.Tier)
				case "boost":
					key = fmt.Sprint(p.AdventureBoost.Multiplier)
				case "adventure":
					key = p.AdventureType
				case "craving":
					key = fmt.Sprint(p.Craving.ItemID)
				default:
					fmt.Fprintf(os.Stderr, "Invalid group field: %s\n", groupKey)
					os.Exit(1)
				}
				groups[key] = append(groups[key], p)
			}

			// Sort group keys
			var keys []string
			for k := range groups {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			// Print per group
			for _, k := range keys {
				items := groups[k]
				fmt.Printf("%s (%d)\n", k, len(items))

				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "ID\tName\tSpecies\tTier\tXP\tAdventure\tItems\tBoost\tEnds")
				for _, p := range items {
					fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%d\t%s\t%d\t%d\t%s\n",
						p.ID, p.Name, p.Species, p.Tier, p.XP,
						p.AdventureType, p.LifetimeItemsFound,
						p.AdventureBoost.Multiplier, common.EpochToISO8601(p.AdventureBoost.EndTime))
				}
				w.Flush()
				fmt.Println()
			}
			return
		}

		// Default table
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tName\tSpecies\tTier\tXP\tAdventure\tItems\tBoost\tEnds")
		for _, p := range pets {
			fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%d\t%s\t%d\t%d\t%s\n",
				p.ID, p.Name, p.Species, p.Tier, p.XP,
				p.AdventureType, p.LifetimeItemsFound,
				p.AdventureBoost.Multiplier, common.EpochToISO8601(p.AdventureBoost.EndTime))
		}
		w.Flush()
	},
}

var offspringCmd = &cobra.Command{
	Use:   "offspring [id]",
	Short: "Fetch offspring",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := common.ParseID(args[0])
		payload := map[string]interface{}{"type": "petOffspring", "id": id}
		data := client.FetchDataOrExit(payload)
		common.PrintJSON(data)
	},
}
