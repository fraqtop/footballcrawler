package competition

import (
	"database/sql"
	"github.com/fraqtop/footballcore/competition"
	"github.com/fraqtop/footballcrawler/connection"
	"github.com/fraqtop/footballcrawler/source"
)

type readRepository struct {
	connection *sql.DB
}

func (r readRepository) All() []competition.Competition {
	var result []competition.Competition

	query, err := r.connection.Query("select id, title, uri_table from competition")
	if err != nil {
		panic("can not fetch data from competition source!")
	}
	defer query.Close()

	for query.Next() {
		var (
			id int
			title string
			uri string
		)

		query.Scan(&id, &title, &uri)
		result = append(result, source.NewCompetitionSource(id, title, uri))
	}

	return result
}

func NewReadRepository() (competition.ReadRepository, error) {
	dbConnection, err := connection.GetInstance()
	if err != nil {
		return nil, err
	}

	return readRepository{
		connection: dbConnection,
	}, nil
}
