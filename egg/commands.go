package egg

import (
	"bcncli/client"
	"bcncli/common"
	"encoding/json"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "egg",
	Short: "Manage eggs",
}

func init() {
	Cmd.AddCommand(infoCmd, ownedCmd, offspringCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info [id]",
	Short: "Fetch egg info",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := common.ParseID(args[0])
		payload := map[string]interface{}{"type": "egg", "id": id}
		data := client.FetchDataOrExit(payload)
		common.PrintJSON(data)
	},
}

var ownedCmd = &cobra.Command{
	Use:   "owned [userId]",
	Short: "List eggs for a user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userId := common.ParseID(args[0])
		payload := map[string]interface{}{"type": "userPetsAndEggs", "id": userId}
		raw := client.FetchDataOrExit(payload)
		var resp map[string]json.RawMessage
		json.Unmarshal(raw, &resp)
		common.PrintJSON(resp["eggs"])
	},
}

var offspringCmd = &cobra.Command{
	Use:   "offspring [id]",
	Short: "Fetch offspring for an egg",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := common.ParseID(args[0])
		payload := map[string]interface{}{"type": "petOffspring", "id": id}
		data := client.FetchDataOrExit(payload)
		common.PrintJSON(data)
	},
}
