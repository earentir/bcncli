package faction

import (
	"fmt"
	"os"

	"bcncli/client"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "faction",
	Short: "Manage factions",
}

func init() {
	Cmd.AddCommand(infoCmd, membersCmd, recruitingCmd, requestsCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info [id]",
	Short: "Fetch faction info",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := client.ParseID(args[0])
		payload := map[string]interface{}{"type": "faction", "id": id}
		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}

var membersCmd = &cobra.Command{
	Use:   "members [id]",
	Short: "List faction members",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := client.ParseID(args[0])
		payload := map[string]interface{}{"type": "factionMembers", "id": id}
		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}

var recruitingCmd = &cobra.Command{
	Use:   "recruiting",
	Short: "List recruiting factions",
	Run: func(cmd *cobra.Command, args []string) {
		payload := map[string]interface{}{"type": "recruitingFactions"}
		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}

var requestsCmd = &cobra.Command{
	Use:   "requests <faction|user> [id]",
	Short: "List join requests by faction or user",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		idTypeArg := args[0]
		id := client.ParseID(args[1])
		var idType string
		if idTypeArg == "faction" {
			idType = "factionId"
		} else if idTypeArg == "user" {
			idType = "bcId"
		} else {
			fmt.Fprintf(os.Stderr, "Invalid type %s, must be 'faction' or 'user'\n", idTypeArg)
			os.Exit(1)
		}
		payload := map[string]interface{}{"type": "factionJoinRequests", "idType": idType, "id": id}
		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}
