package profile

import (
	"bcncli/client"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage user profiles",
}

func init() {
	Cmd.AddCommand(infoCmd, userCmd, inventoryCmd, statsCmd, trophiesCmd, flatinventoryCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info [id]",
	Short: "Fetch profile info",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := client.ParseID(args[0])
		payload := map[string]interface{}{"type": "profile", "id": id}
		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}

var userCmd = &cobra.Command{
	Use:   "user [id]",
	Short: "Fetch user details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := client.ParseID(args[0])
		payload := map[string]interface{}{"type": "user", "id": id}
		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}

var inventoryCmd = &cobra.Command{
	Use:   "inventory [id]",
	Short: "Fetch inventory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := client.ParseID(args[0])
		payload := map[string]interface{}{"type": "inventory", "id": id}
		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}

var flatinventoryCmd = &cobra.Command{
	Use:   "flatinventory [id]",
	Short: "Fetch flat inventory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := client.ParseID(args[0])
		payload := map[string]interface{}{"type": "flatInventory", "id": id}
		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats [id]",
	Short: "Fetch stats",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := client.ParseID(args[0])
		payload := map[string]interface{}{"type": "stats", "id": id}
		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}

var trophiesCmd = &cobra.Command{
	Use:   "trophies [id]",
	Short: "Fetch trophies",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := client.ParseID(args[0])
		payload := map[string]interface{}{"type": "trophies", "id": id}
		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}
