package shared

import (
	"os"
	"strings"
)

type Config struct {
	Port        string
	DatabaseURL string
	CorsOrigins []string
	LogPayloads bool
}

func LoadConfig() Config {
	return Config{
		Port:        os.Getenv("PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		CorsOrigins: strings.Split(os.Getenv("CORS_ORIGINS"), ","),
		LogPayloads: os.Getenv("LOG_PAYLOADS") == "true",
	}
}
