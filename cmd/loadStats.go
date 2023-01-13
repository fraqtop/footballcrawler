package cmd

import (
	"fmt"
	corecompetitions "github.com/fraqtop/footballcore/competition"
	corestats "github.com/fraqtop/footballcore/stats"
	"github.com/fraqtop/footballcrawler/internal/repository/competition"
	"github.com/fraqtop/footballcrawler/internal/repository/stats"
	"github.com/spf13/cobra"
	"os"
)

var loadStatsCmd = &cobra.Command{
	Use:   "loadStats",
	Short: "Load statistics of predefined competitions",
	Run: func(cmd *cobra.Command, args []string) {
		competitionRepository, err := competition.NewReadRepository()
		if err != nil {
			fmt.Println(ErrConnectToStorage)
			os.Exit(1)
		}
		statsReadRepository := stats.NewReadRepository()
		statsWriteRepository, err := stats.NewWriteRepository()
		if err != nil {
			panic(err)
		}

		statsChannel := make(chan []corestats.Stats)
		threadsCount := 0
		for _, competitionEntity := range competitionRepository.All() {
			go func(outputChannel chan []corestats.Stats, competitionEntity corecompetitions.Competition, repository corestats.ReadRepository) {
				statsChannel <- statsReadRepository.ByCompetition(competitionEntity)
			}(statsChannel, competitionEntity, statsReadRepository)
			threadsCount++
		}

		if threadsCount == 0 {
			fmt.Println("no competitions loaded, try to seed them first and then try again")
			os.Exit(1)
		}

		var foundStats []corestats.Stats
		for output := range statsChannel {
			foundStats = append(foundStats, output...)
			threadsCount--
			if threadsCount == 0 {
				close(statsChannel)
			}
		}

		if err = statsWriteRepository.BatchUpdate(foundStats); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loadStatsCmd)
}
