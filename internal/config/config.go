package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Postgres
	LogLevel      string `env:"LOG_LEVEL"`
	ApiSecret     string `env:"API_SECRET"`
	ServerAddress string `env:"SERVER_ADDR"`
	ServerPort    string `env:"SERVER_PORT"`
}

type Postgres struct {
	Url          string   `env:"DB_URL"`
	Repository   string   `env:"REPOSITORY"`
	Version      string   `env:"VERSION"`
	CfgVariables []string `env:"REPO_CONFIG"`
	Name         string   `env:"POSTGRES_DB"`
	User         string   `env:"POSTGRES_USER"`
	Password     string   `env:"POSTGRES_PASSWORD"`
}

func NewConfig() *Config {
	cfg := Config{}

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal().Msgf("error loading .env: %v", err)
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Msgf("error parsing .env: %v", err)
	}

	return &cfg
}
