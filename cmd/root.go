package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "footballcrawler",
	Short: "Fill env database with football statistics",
	Long: `This tool fetches competitions, teams, statistics etc from the source and upload them 
to given database in current environment`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
