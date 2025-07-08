package pet

import (
	"encoding/json"

	"bcncli/client"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "pet",
	Short: "Manage pets",
}

func init() {
	Cmd.AddCommand(infoCmd, ownedCmd, offspringCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info [id]",
	Short: "Fetch pet info",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := client.ParseID(args[0])
		payload := map[string]interface{}{"type": "pet", "id": id}
		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}

var ownedCmd = &cobra.Command{
	Use:   "owned [userId]",
	Short: "List pets for a user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userId := client.ParseID(args[0])
		payload := map[string]interface{}{"type": "userPetsAndEggs", "id": userId}
		raw := client.FetchDataOrExit(payload)
		var resp map[string]json.RawMessage
		json.Unmarshal(raw, &resp)
		client.PrintJSON(resp["pets"])
	},
}

var offspringCmd = &cobra.Command{
	Use:   "offspring [id]",
	Short: "Fetch offspring",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := client.ParseID(args[0])
		payload := map[string]interface{}{"type": "petOffspring", "id": id}
		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}
