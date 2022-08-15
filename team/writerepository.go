package team

import (
	"database/sql"
	"fmt"
	"github.com/fraqtop/footballcore/team"
	"github.com/fraqtop/footballcrawler/connection"
	"strings"
)

type writeRepository struct {
	connection *sql.DB
}

func NewWriteRepository() (team.WriteRepository, error) {
	dbConnection, err := connection.GetInstance()
	if err != nil {
		return nil, err
	}
	return writeRepository{connection: dbConnection}, nil
}

func (w writeRepository) BatchUpdate(teams []team.Team) error {
	var insertQueryParts []string

	for _, entity := range teams {
		currentQueryPart := fmt.Sprintf("('%s', '%s')", entity.TitleShort(), entity.TitleFull())
		insertQueryParts = append(insertQueryParts, currentQueryPart)
	}

	queryPattern := "insert into team (title_short, title_full) " +
		"values %s " +
		"on conflict (title_short, title_full) do update " +
		"set title_full = excluded.title_full " +
		"returning id"

	queryToExecute := fmt.Sprintf(queryPattern, strings.Join(insertQueryParts, ","))
	rows, err := w.connection.Query(queryToExecute)

	if err != nil {
		return err
	}
	defer rows.Close()

	var id int
	i := 0
	for rows.Next() {
		_ = rows.Scan(&id)
		teams[i].SetId(id)
		i++
	}

	return nil
}
