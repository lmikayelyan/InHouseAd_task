package pkg

import (
	"InHouseAd/internal/config"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Session interface {
	InitPgSession(ctx context.Context) (*pgxpool.Pool, error)
}

type session struct {
	cfg *config.Postgres
}

func NewPgSession(cfg *config.Postgres) Session {
	return &session{cfg: cfg}
}

func (s *session) InitPgSession(ctx context.Context) (*pgxpool.Pool, error) {
	connConfig, err := pgxpool.ParseConfig(s.cfg.Url)
	if err != nil {
		log.Panic(err)
	}

	pool, err := pgxpool.ConnectConfig(ctx, connConfig)
	if err != nil {
		//log.Panic(err)
	}

	return pool, nil
}
