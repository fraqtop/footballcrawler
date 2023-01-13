package team

import (
	"errors"
	"fmt"
	"github.com/fraqtop/footballcore/competition"
	"github.com/fraqtop/footballcore/team"
	"github.com/fraqtop/footballcrawler/internal/source"
	"github.com/gocolly/colly/v2"
	"os"
)

type readRepository struct {
	sourceReadRepository source.ReadRepository
}

func (trr readRepository) ByCompetition(competition competition.Competition) []team.Team {
	sourceEntity, ok := competition.(source.CompetitionSource)
	if ok {
		return trr.parseSource(sourceEntity)
	}
	sourceEntity, err := trr.findSource(competition)
	if err != nil {
		return []team.Team{}
	}

	return trr.parseSource(sourceEntity)
}

func (trr readRepository) findSource(requestedCompetition competition.Competition) (source.CompetitionSource, error) {
	for _, competitionSource := range trr.sourceReadRepository.Competitions() {
		if competitionSource.Id() == requestedCompetition.Id() {
			return competitionSource, nil
		}
	}

	return source.CompetitionSource{}, errors.New("can't find suitable source")
}

func (trr readRepository) parseSource(competitionSource source.CompetitionSource) []team.Team {
	collector := colly.NewCollector()

	var teamLabels []string
	collector.OnHTML("a[data-clubhouse-uid]", func(element *colly.HTMLElement) {
		text := element.Text

		if len(text) > 0 {
			teamLabels = append(teamLabels, text)
		}
	})

	if err := collector.Visit(competitionSource.Uri()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var (
		teams             []team.Team
		currentShortTitle string
		currentFullTitle  string
	)

	for i := 0; i < len(teamLabels); i++ {
		if i%2 == 0 {
			currentShortTitle = teamLabels[i]
		} else {
			currentFullTitle = teamLabels[i]
			teams = append(teams, team.New(0, currentFullTitle, currentShortTitle))
		}
	}

	return teams
}

func NewReadRepository() team.ReadRepository {
	return readRepository{
		sourceReadRepository: source.Instance(),
	}
}
