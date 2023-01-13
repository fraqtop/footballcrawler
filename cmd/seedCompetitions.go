package cmd

import (
	"fmt"
	"github.com/fraqtop/footballcrawler/internal/repository/competition"
	"github.com/fraqtop/footballcrawler/internal/source"
	"os"

	"github.com/spf13/cobra"
)

const ErrConnectToStorage  = "sorry, can't connect to storage"

var seedCompetitionsCmd = &cobra.Command{
	Use:   "seedCompetitions",
	Short: "Loads predefined competitions to database",
	Run: func(cmd *cobra.Command, args []string) {
		readRepository := source.Instance()
		writeRepository, err := competition.NewWriteRepository()
		if err != nil {
			fmt.Println(ErrConnectToStorage)
			os.Exit(1)
		}

		for _, competitionEntity := range readRepository.Competitions() {
			if err := writeRepository.Save(competitionEntity); err != nil {
				fmt.Println("can't save competition", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(seedCompetitionsCmd)
}
