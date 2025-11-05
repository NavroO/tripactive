package shared

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func ConnectDB() (*sql.DB, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal().Msg("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Info().Msg("Connected to PostgreSQL")
	return db, nil
}
