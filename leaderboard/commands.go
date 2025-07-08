// File: leaderboard/commands.go
package leaderboard

import (
	"fmt"
	"os"

	"bcncli/client"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "leaderboard",
	Short: "View various leaderboards",
}

func init() {
	Cmd.AddCommand(userCmd, factionCmd, petsCmd)

	userCmd.Flags().StringP("lbType", "t", "", "Leaderboard type: rank, questLevel, stat, or item")
	userCmd.Flags().StringP("stat", "s", "", "Statistic name (required if --lbType=stat)")
	userCmd.Flags().IntP("itemId", "i", 0, "Item ID (required if --lbType=item)")
	userCmd.Flags().IntP("page", "p", 1, "Page number")
	userCmd.MarkFlagRequired("lbType")
	userCmd.MarkFlagRequired("page")

	factionCmd.Flags().StringP("stat", "s", "", "fpDepositedMonthly or fpDepositedTotal")
	factionCmd.Flags().IntP("page", "p", 1, "Page number")
	factionCmd.MarkFlagRequired("stat")
	factionCmd.MarkFlagRequired("page")

	petsCmd.Flags().IntP("page", "p", 1, "Page number")
	petsCmd.MarkFlagRequired("page")
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Show user leaderboard",
	Run: func(cmd *cobra.Command, args []string) {
		lbType, _ := cmd.Flags().GetString("lbType")
		stat, _ := cmd.Flags().GetString("stat")
		itemId, _ := cmd.Flags().GetInt("itemId")
		page, _ := cmd.Flags().GetInt("page")

		// Validate conditional flags
		if lbType == "stat" && stat == "" {
			fmt.Fprintln(os.Stderr, "Error: --stat is required when --lbType=stat")
			os.Exit(1)
		}
		if lbType == "item" && itemId == 0 {
			fmt.Fprintln(os.Stderr, "Error: --itemId is required when --lbType=item")
			os.Exit(1)
		}

		payload := map[string]interface{}{
			"type":   "userLeaderboard",
			"lbType": lbType,
			"page":   page,
		}
		if stat != "" {
			payload["stat"] = stat
		}
		if itemId != 0 {
			payload["itemId"] = itemId
		}

		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}

var factionCmd = &cobra.Command{
	Use:   "faction",
	Short: "Show faction leaderboard",
	Run: func(cmd *cobra.Command, args []string) {
		stat, _ := cmd.Flags().GetString("stat")
		page, _ := cmd.Flags().GetInt("page")

		payload := map[string]interface{}{
			"type": "factionLeaderboard",
			"stat": stat,
			"page": page,
		}

		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}

var petsCmd = &cobra.Command{
	Use:   "pets",
	Short: "Show pets leaderboard",
	Run: func(cmd *cobra.Command, args []string) {
		page, _ := cmd.Flags().GetInt("page")

		payload := map[string]interface{}{
			"type": "petsLeaderboard",
			"page": page,
		}

		data := client.FetchDataOrExit(payload)
		client.PrintJSON(data)
	},
}
