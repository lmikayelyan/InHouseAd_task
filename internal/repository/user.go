package repository

import (
	"InHouseAd/internal/model"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User interface {
	Create(ctx context.Context, user model.User) error
	GetHash(ctx context.Context, username string) (string, error)
	GetID(ctx context.Context, username string) (uint, error)
}

type user struct {
	pool *pgxpool.Pool
}

func NewUser(pool *pgxpool.Pool) User {
	return &user{pool: pool}
}

func (u *user) Create(ctx context.Context, user model.User) error {
	queryStr := "insert into users(username, password, e_mail, phone_number) VALUES($1, $2, $3, $4)"
	_, err := u.pool.Exec(ctx, queryStr,
		user.Username,
		user.Password,
		user.EMail,
		user.PhoneNumber,
	)

	return err
}

func (u *user) GetHash(ctx context.Context, username string) (string, error) {
	queryStr := "select password from users where username=$1"
	var pass string
	err := u.pool.QueryRow(ctx, queryStr, username).Scan(&pass)
	return pass, err
}

func (u *user) GetID(ctx context.Context, username string) (uint, error) {
	var id uint
	err := u.pool.QueryRow(ctx, "select id from users where username=$1", username).Scan(&id)

	return id, err
}
