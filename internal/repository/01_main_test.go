package repository_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"
	"github.com/rs/zerolog/log"
)

var TestPool *pgxpool.Pool

func TestMain(m *testing.M) {
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal().Msgf("Could not connect to docker: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresResource := initializePostgres(ctx, dockerPool, newPostgresConfig())
	code := m.Run()

	purgeResources(dockerPool, postgresResource)

	os.Exit(code)
}

func initializePostgres(ctx context.Context, dockerPool *dockertest.Pool, cfg *postgresConfig) *dockertest.Resource {
	resource, err := dockerPool.Run(cfg.Repository, cfg.Version, cfg.EnvVariables)
	if err != nil {
		log.Fatal().Msgf("Could not start resource: %s", err)
	}

	var dbHostAndPort string

	err = dockerPool.Retry(func() error {
		var dbHost string

		gitlabCIHost := os.Getenv("DATABASE_HOST")

		if gitlabCIHost != "" {
			dbHost = gitlabCIHost
		} else {
			dbHost = "localhost"
		}

		port := resource.GetPort(cfg.PortID)
		dbHostAndPort = fmt.Sprintf("%s:%s", dbHost, port)

		dsn := cfg.getConnectionString(dbHostAndPort)

		TestPool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return fmt.Errorf("connect: %v", err)
		}

		if err = TestPool.Ping(ctx); err != nil {
			return fmt.Errorf("ping: %v", err)
		}

		return nil
	})
	if err != nil {
		log.Fatal().Msgf("Could not connect to database: %s", err)
	}
	log.Info().Msg(strings.Join(cfg.getFlywayMigrationArgs(dbHostAndPort), " "))
	cmd := exec.Command("flyway", cfg.getFlywayMigrationArgs(dbHostAndPort)...)

	err = cmd.Run()
	if err != nil {
		log.Fatal().Msgf("There are errors in migrations: %v", err)
	}
	return resource
}

func purgeResources(dockerPool *dockertest.Pool, resources ...*dockertest.Resource) {
	for i := range resources {
		if err := dockerPool.Purge(resources[i]); err != nil {
			log.Fatal().Msgf("Could not purge resource: %s", err)
		}
		err := resources[i].Expire(1)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	}
}

type postgresConfig struct {
	Repository   string
	Version      string
	EnvVariables []string
	PortID       string
	DB           string
}

func newPostgresConfig() *postgresConfig {
	return &postgresConfig{
		Repository: "postgres",
		Version:    "14.1-alpine",
		EnvVariables: []string{
			"POSTGRES_PASSWORD=password123",
			"POSTGRES_DB=referral",
		},
		PortID: "5432/tcp",
		DB:     "referral",
	}
}

func (p *postgresConfig) getConnectionString(dbHostAndPort string) string {
	return fmt.Sprintf("postgresql://postgres:password123@%v/%s?sslmode=disable", dbHostAndPort, p.DB)
}

func (p *postgresConfig) getFlywayMigrationArgs(dbHostAndPort string) []string {
	return []string{
		"-user=postgres",
		"-password=password123",
		"-locations=filesystem:../../migration",
		fmt.Sprintf("-url=jdbc:postgresql://%v/%s", dbHostAndPort, p.DB),
		"migrate",
	}
}
