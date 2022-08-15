package stats

import (
	"database/sql"
	"fmt"
	corecompetition "github.com/fraqtop/footballcore/competition"
	"github.com/fraqtop/footballcore/stats"
	coreteam "github.com/fraqtop/footballcore/team"
	"github.com/fraqtop/footballcrawler/competition"
	"github.com/fraqtop/footballcrawler/connection"
	"github.com/fraqtop/footballcrawler/team"
	"strings"
)

type writeRepository struct {
	connection                 *sql.DB
	competitionWriteRepository corecompetition.WriteRepository
	teamWriteRepository        coreteam.WriteRepository
}

func (w writeRepository) BatchUpdate(stats []stats.Stats) error {
	if err := w.syncTeams(stats); err != nil {
		return err
	}
	var insertQueryParts []string
	distinctCompetitions := make(map[int]corecompetition.Competition)
	for _, currentStats := range stats {
		distinctCompetitions[currentStats.Competition().Id()] = currentStats.Competition()
		insertQueryParts = append(
			insertQueryParts,
			fmt.Sprintf(
				"(%d, %d, '%s', %d, %d, %d, %d, %d, %d, %d)",
				currentStats.Team().Id(),
				currentStats.Competition().Id(),
				currentStats.Season(),
				currentStats.Games(),
				currentStats.Points(),
				currentStats.Wins(),
				currentStats.Draws(),
				currentStats.Losses(),
				currentStats.Scored(),
				currentStats.Passed(),
			),
		)
	}

	if err := w.syncCompetitions(distinctCompetitions); err != nil {
		return err
	}

	sqlPattern := "insert into stats (team_id, competition_id, season, games, points, wins, draws, losses, scored, passed)" +
		"values %s " +
		"on conflict(team_id, competition_id) do update set " +
		"season = excluded.season," +
		"games = excluded.games," +
		"points = excluded.points," +
		"wins = excluded.wins," +
		"draws = excluded.draws," +
		"losses = excluded.losses," +
		"scored = excluded.scored," +
		"passed = excluded.passed;"

	queryToExecute := fmt.Sprintf(sqlPattern, strings.Join(insertQueryParts, ","))
	_, err := w.connection.Exec(queryToExecute)

	return err
}

func (w writeRepository) syncTeams(stats []stats.Stats) error {
	distinctTeams := make(map[string]coreteam.Team)
	for _, currentStats := range stats {
		distinctTeams[currentStats.Team().TitleShort() + currentStats.Team().TitleFull()] = currentStats.Team()
	}

	var teams []coreteam.Team
	for _, currentTeam := range distinctTeams {
		teams = append(teams, currentTeam)
	}

	return w.teamWriteRepository.BatchUpdate(teams)
}

func (w writeRepository) syncCompetitions(competitions map[int]corecompetition.Competition) error {
	for _, currentCompetition := range competitions {
		if err := w.competitionWriteRepository.Save(currentCompetition); err != nil {
			return err
		}
	}

	return nil
}

func NewWriteRepository() (stats.WriteRepository, error) {
	dbConnection, err := connection.GetInstance()
	if err != nil {
		return nil, err
	}
	teamWriteRepository, err := team.NewWriteRepository()
	if err != nil {
		return nil, err
	}
	competitionRepository, err := competition.NewWriteRepository()
	if err != nil {
		return nil, err
	}

	return writeRepository{
			connection:                 dbConnection,
			competitionWriteRepository: competitionRepository,
			teamWriteRepository:        teamWriteRepository,
		},
		nil
}
