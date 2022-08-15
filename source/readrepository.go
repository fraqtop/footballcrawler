package source

import (
	"github.com/fraqtop/footballcore/competition"
	"os"
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
	sources := []CompetitionSource {
		NewCompetitionSource(1, "Spanish La Liga", os.Getenv("CRAWLER_HOST") + "/soccer/standings/_/league/esp.1/season/2021"),
		NewCompetitionSource(2, "English Premier League", os.Getenv("CRAWLER_HOST") + "/soccer/standings/_/league/eng.1/season/2021"),
		NewCompetitionSource(3, "German Bundesliga", os.Getenv("CRAWLER_HOST") + "/soccer/standings/_/league/ger.1/season/2021"),
		NewCompetitionSource(4, "Italian Serie A", os.Getenv("CRAWLER_HOST") + "/soccer/standings/_/league/ita.1/season/2021"),
		NewCompetitionSource(5, "French Ligue 1", os.Getenv("CRAWLER_HOST") + "/soccer/standings/_/league/fra.1/season/2021"),
		NewCompetitionSource(6, "Russian Premier League", os.Getenv("CRAWLER_HOST") + "/soccer/standings/_/league/rus.1/season/2021"),

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
