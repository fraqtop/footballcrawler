package source

import (
	"fmt"
	"os"
	"time"

	"github.com/fraqtop/footballcore/competition"
)

var readRepositoryInstance *readRepository = nil

type ReadRepository interface {
	Competitions() []CompetitionSource
}

type CompetitionSource struct {
	uri              string
	innerCompetition competition.Competition
}

type readRepository struct {
	sources []CompetitionSource
}

func (r readRepository) Competitions() []CompetitionSource {
	return r.sources
}

func (c CompetitionSource) Id() int {
	return c.innerCompetition.Id()
}

func (c CompetitionSource) Title() string {
	return c.innerCompetition.Title()
}

func (c CompetitionSource) Uri() string {
	return c.uri
}

func Instance() ReadRepository {
	if readRepositoryInstance == nil {
		initReadRepositoryInstance()
	}

	return readRepositoryInstance
}

func initReadRepositoryInstance() {
	year, month, _ := time.Now().Date()
	//season started in prev year
	if month < 9 {
		year--
	}
	fullNames := []string{
		"Spanish La Liga",
		"English Premier League",
		"German Bundesliga",
		"Italian Serie A",
		"French Ligue 1",
		"Russian Premier League",
	}
	uriNames := []string{
		"esp.1",
		"eng.1",
		"ger.1",
		"ita.1",
		"fra.1",
		"rus.1",
	}

	sources := make([]CompetitionSource, 0, len(fullNames))
	for i, currFullName := range fullNames {
		sources = append(sources, NewCompetitionSource(
			i+1,
			currFullName,
			fmt.Sprintf(
				"%s/soccer/standings/_/league/%s/season/%d",
				os.Getenv("CRAWLER_HOST"),
				uriNames[i],
				year,
			)),
		)
	}

	readRepositoryInstance = &readRepository{
		sources: sources,
	}
}

func NewCompetitionSource(id int, title, uri string) CompetitionSource {
	newCompetition := competition.New(id, title)

	return CompetitionSource{
		uri:              uri,
		innerCompetition: newCompetition,
	}
}
