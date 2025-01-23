package cmd

import (
	"fmt"
	"os"

	"github.com/fraqtop/footballcrawler/internal/repository/competition"
	"github.com/fraqtop/footballcrawler/internal/source"

	"github.com/spf13/cobra"
)

const ErrConnectToStorage = "sorry, can't connect to storage"

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
				fmt.Printf("can't save competition %s cause %s\n", competitionEntity.Title(), err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(seedCompetitionsCmd)
}
