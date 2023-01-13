package connection

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

var (
	connectionInstance                *sql.DB = nil
	errDatabaseConfigurationIncorrect         = errors.New("can't connect to database")
)

func GetInstance() (*sql.DB, error) {
	if connectionInstance == nil {
		var err error
		connectionInstance, err = sql.Open(
			"postgres",
			fmt.Sprintf(
				"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
				os.Getenv("POSTGRES_USER"),
				os.Getenv("POSTGRES_PASSWORD"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_EXPOSE_PORT"),
				os.Getenv("POSTGRES_DB"),
				"disable",
			),
		)

		if err != nil {
			return nil, errDatabaseConfigurationIncorrect
		}
	}

	return connectionInstance, nil
}
