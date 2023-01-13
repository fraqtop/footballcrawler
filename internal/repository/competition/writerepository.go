package competition

import (
	"database/sql"
	"github.com/fraqtop/footballcore/competition"
	"github.com/fraqtop/footballcrawler/internal/connection"
	"github.com/fraqtop/footballcrawler/internal/source"
	"os"
)

type writeRepository struct {
	connection *sql.DB
}

func (w writeRepository) Save(competition competition.Competition) error {
	var err error
	if sourceEntity, ok := competition.(source.CompetitionSource); ok {
		_, err = w.connection.Exec("insert into competition (id, title, uri_table) "+
			"values ($1, $2, $3) "+
			"on conflict(id) do update "+
			"set title = excluded.title, uri_table = excluded.uri_table", sourceEntity.Id(), sourceEntity.Title(), sourceEntity.Uri())
		if err != nil {
			return err
		}
	} else {
		_, err = w.connection.Exec("insert into competition (id, title, uri_table) "+
			"values ($1, $2, $3) "+
			"on conflict(id) do update "+
			"set title = excluded.title", competition.Id(), competition.Title(), os.Getenv("CRAWLER_HOST"))
		if err != nil {
			return err
		}
	}

	return err
}

func NewWriteRepository() (competition.WriteRepository, error) {
	dbConnection, err := connection.GetInstance()
	if err != nil {
		return nil, err
	}

	return writeRepository{connection: dbConnection}, nil
}
