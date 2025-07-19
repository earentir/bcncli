// File: search/commands.go
package search

import (
	"bcncli/client"
	"bcncli/common"

	"github.com/spf13/cobra"
)

// Cmd is the root command for search operations
var Cmd = &cobra.Command{
	Use:   "search",
	Short: "Search for users, factions, or pets",
}

func init() {
	Cmd.AddCommand(userCmd, factionCmd, petCmd)

	// pet command flags
	petCmd.Flags().StringP("skin", "s", "any skin", "Skin filter: specific, 'no skin', or 'any skin'")
	petCmd.Flags().StringP("aura", "a", "any aura", "Aura filter: specific, 'no aura', or 'any aura'")
	petCmd.Flags().StringP("species", "c", "any species", "Species filter: specific or 'any species'")
	petCmd.Flags().StringP("name", "n", "", "Name query filter (rawNameQuery)")
}

// userCmd searches for users by name
var userCmd = &cobra.Command{
	Use:   "user [query]",
	Short: "Search users by name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		payload := map[string]interface{}{
			"type":  "searchUsers",
			"query": query,
		}
		data := client.FetchDataOrExit(payload)
		common.PrintJSON(data)
	},
}

// factionCmd searches for factions by name
var factionCmd = &cobra.Command{
	Use:   "faction [query]",
	Short: "Search factions by name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		payload := map[string]interface{}{
			"type":  "searchFactions",
			"query": query,
		}
		data := client.FetchDataOrExit(payload)
		common.PrintJSON(data)
	},
}

// petCmd searches for pets with filters
var petCmd = &cobra.Command{
	Use:   "pet",
	Short: "Search pets by properties",
	Run: func(cmd *cobra.Command, args []string) {
		skin, _ := cmd.Flags().GetString("skin")
		aura, _ := cmd.Flags().GetString("aura")
		species, _ := cmd.Flags().GetString("species")
		name, _ := cmd.Flags().GetString("name")

		payload := map[string]interface{}{
			"type":         "searchPets",
			"skin":         skin,
			"aura":         aura,
			"species":      species,
			"rawNameQuery": name,
		}
		data := client.FetchDataOrExit(payload)
		common.PrintJSON(data)
	},
}
