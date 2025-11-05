package shared

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var LogPayloads bool

func SetupLogger() {
	zerolog.TimeFieldFormat = time.RFC3339

	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
	}).With().Timestamp().Logger()

	LogPayloads = os.Getenv("LOG_PAYLOADS") == "true"
}
