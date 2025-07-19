// File: logs/commands.go
package logs

import (
	"fmt"
	"os"

	"bcncli/client"
	"bcncli/common"

	"github.com/spf13/cobra"
)

// Cmd is the root command for log-related operations
var Cmd = &cobra.Command{
	Use:   "logs",
	Short: "Retrieve various logs",
}

func init() {
	Cmd.AddCommand(bcidCmd, idtypeCmd, logtypeCmd, inputsCmd)

	// page flag for commands that support it
	bcidCmd.Flags().IntP("page", "p", 1, "Page number")
	idtypeCmd.Flags().IntP("page", "p", 1, "Page number")
	logtypeCmd.Flags().IntP("page", "p", 1, "Page number")
}

// bcidCmd fetches logs by a user's BCID; bcId is positional, page is a flag
var bcidCmd = &cobra.Command{
	Use:   "bcid [bcId]",
	Short: "List logs for a user by BCID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bcId := common.ParseID(args[0])
		page, _ := cmd.Flags().GetInt("page")

		payload := map[string]interface{}{
			"type": "richLogsByBcId",
			"id":   bcId,
			"page": page,
		}
		data := client.FetchDataOrExit(payload)
		common.PrintJSON(data)
	},
}

// idtypeCmd fetches logs by ID type (faction or item); idType and id are positional, page is a flag
var idtypeCmd = &cobra.Command{
	Use:   "idtype [faction|item] [id]",
	Short: "List logs by ID type (factionId or itemId)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		t := args[0]
		var idType string
		switch t {
		case "faction":
			idType = "factionId"
		case "item":
			idType = "itemId"
		default:
			fmt.Fprintf(os.Stderr, "Invalid ID type '%s', must be 'faction' or 'item'\n", t)
			os.Exit(1)
		}
		id := common.ParseID(args[1])
		page, _ := cmd.Flags().GetInt("page")

		payload := map[string]interface{}{
			"type":   "richLogsByIdType",
			"idType": idType,
			"id":     id,
			"page":   page,
		}
		data := client.FetchDataOrExit(payload)
		common.PrintJSON(data)
	},
}

// logtypeCmd fetches logs by log type; logType is positional, page is a flag
var logtypeCmd = &cobra.Command{
	Use:   "logtype [logType]",
	Short: "List logs by log type",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logType := args[0]
		page, _ := cmd.Flags().GetInt("page")

		payload := map[string]interface{}{
			"type":    "richLogsByLogType",
			"logType": logType,
			"page":    page,
		}
		data := client.FetchDataOrExit(payload)
		common.PrintJSON(data)
	},
}

// inputsCmd fetches daily user inputs for a given date; bcId and date are positional
var inputsCmd = &cobra.Command{
	Use:   "inputs [bcId] [date]",
	Short: "List daily inputs for user on a date",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		bcId := common.ParseID(args[0])
		date := args[1]

		payload := map[string]interface{}{
			"type": "dailyUserInputs",
			"id":   bcId,
			"date": date,
		}
		data := client.FetchDataOrExit(payload)
		common.PrintJSON(data)
	},
}
