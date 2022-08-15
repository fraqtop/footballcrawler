package main

import (
	"flag"
	"fmt"
	corecompetitions "github.com/fraqtop/footballcore/competition"
	corestats "github.com/fraqtop/footballcore/stats"
	"github.com/fraqtop/footballcrawler/competition"
	"github.com/fraqtop/footballcrawler/source"
	"github.com/fraqtop/footballcrawler/stats"
	"github.com/joho/godotenv"
	"os"
)

const (
	ModeLoadingStats     = "load-stats"
	ModeSeedCompetitions = "seed-comps"
	ErrConnectToStorage  = "sorry, can't connect to storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	currentAction := flag.String("mode", ModeLoadingStats, "crawler execution mode")
	flag.Parse()

	switch *currentAction {
	case ModeLoadingStats:
		loadStats()
	case ModeSeedCompetitions:
		seedCompetitions()
	default:
		panic("your action command argument is invalid")
	}
}

func loadStats() {
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

	statsChannel := make (chan []corestats.Stats)
	threadsCount := 0
	for _, competitionEntity := range competitionRepository.All() {
		go func(outputChannel chan []corestats.Stats, competitionEntity corecompetitions.Competition, repository corestats.ReadRepository) {
			statsChannel <- statsReadRepository.ByCompetition(competitionEntity)
		}(statsChannel, competitionEntity, statsReadRepository)
		threadsCount++
	}

	if threadsCount == 0 {
		fmt.Println("no competitions loaded, try to run with -mode seed-comps")
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
}

func seedCompetitions() {
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
}
