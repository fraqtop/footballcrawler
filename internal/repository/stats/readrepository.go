package stats

import (
	"fmt"
	"github.com/fraqtop/footballcore/competition"
	"github.com/fraqtop/footballcore/stats"
	"github.com/fraqtop/footballcore/team"
	crawlerTeam "github.com/fraqtop/footballcrawler/internal/repository/team"
	"github.com/fraqtop/footballcrawler/internal/source"
	"github.com/gocolly/colly/v2"
	"os"
	"strconv"
	"time"
)

const teamStatsCount = 7

type readRepository struct {
	sourceRepository source.ReadRepository
	teamRepository   team.ReadRepository
}

func NewReadRepository() stats.ReadRepository {
	return readRepository{
		sourceRepository: source.Instance(),
		teamRepository:   crawlerTeam.NewReadRepository(),
	}
}

func (srr readRepository) ByCompetition(competition competition.Competition) []stats.Stats {
	var uri string
	if sourceEntity, ok := competition.(source.CompetitionSource); ok {
		uri = sourceEntity.Uri()
	} else {
		uri = srr.findSourceLink(competition)
	}

	var result []stats.Stats
	if len(uri) == 0 {
		return result
	}

	statsValues := srr.fetchStatsValues(uri)
	season := srr.fetchSeason(uri)
	teams := srr.teamRepository.ByCompetition(competition)

	for index, currentTeam := range teams {
		offset := index * teamStatsCount
		teamStats := statsValues[offset : offset+teamStatsCount]
		result = append(
			result,
			stats.New(
				currentTeam,
				competition,
				season,
				teamStats[0],
				teamStats[6],
				teamStats[1],
				teamStats[2],
				teamStats[3],
				teamStats[4],
				teamStats[5],
			),
		)
	}

	return result
}

func (srr readRepository) findSourceLink(requestedCompetition competition.Competition) string {
	for _, competitionSource := range srr.sourceRepository.Competitions() {
		if competitionSource.Id() == requestedCompetition.Id() {
			return competitionSource.Uri()
		}
	}

	return ""
}

func (srr readRepository) fetchStatsValues(url string) []int {
	collector := colly.NewCollector()

	var result []int
	collector.OnHTML("span[class=stat-cell]", func(element *colly.HTMLElement) {
		value, err := strconv.Atoi(element.Text)
		if err != nil {
			fmt.Println("something wrong with your xpath, got " + element.Text)
			os.Exit(1)
		}
		result = append(result, value)
	})

	if err := collector.Visit(url); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return result
}

func (srr readRepository) fetchSeason(uri string) string {
	seasonEndSubstring := uri[len(uri)-4:]
	seasonEndYear, err := strconv.Atoi(seasonEndSubstring)
	if err != nil {
		seasonEndYear = time.Now().Year()
	}

	return strconv.Itoa(seasonEndYear) + "-" + strconv.Itoa(seasonEndYear+1)
}
