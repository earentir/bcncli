package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"bcncli/egg"
	"bcncli/faction"
	"bcncli/market"
	"bcncli/pet"
	"bcncli/profile"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "bcncli",
		Short: "BCN CLI interacts with the bconomy API",
	}

	// Global API key flag & config
	rootCmd.PersistentFlags().String("apikey", "", "BConomy API key (flag, config file, or env var BCONOMYAPI)")
	viper.BindPFlag("apikey", rootCmd.PersistentFlags().Lookup("apikey"))
	viper.BindEnv("apikey", "BCONOMYAPI")

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	cobra.OnInitialize(func() {
		if err := viper.ReadInConfig(); err == nil {
			fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
		}
	})

	// Register commands
	rootCmd.AddCommand(pet.Cmd)
	rootCmd.AddCommand(egg.Cmd)
	rootCmd.AddCommand(profile.Cmd)
	rootCmd.AddCommand(faction.Cmd)
	rootCmd.AddCommand(market.Cmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
