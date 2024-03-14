package service

import (
	"InHouseAd/internal/model"
	"InHouseAd/internal/repository"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type User interface {
	Register(ctx context.Context, reg model.Register) error
	Login(ctx context.Context, log model.Login) (map[string]string, error)
	LoginCheck(ctx context.Context, user model.User) (Token, uint, error)
}

type user struct {
	repo repository.User
}

func NewUser(repo repository.User) User {
	return &user{repo: repo}
}

func (u *user) Register(ctx context.Context, reg model.Register) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(reg.Password), 14)
	if err != nil {
		return fmt.Errorf("service.user.Register.GenerateFromPassword: %v", err)
	}

	userInstance := model.User{
		Username:    reg.Username,
		Password:    string(hashedPass),
		EMail:       reg.EMail,
		PhoneNumber: reg.PhoneNumber,
	}

	if err = u.repo.Create(ctx, userInstance); err != nil {
		return fmt.Errorf("service.user.Register.Create: %v", err)
	}

	return nil
}

func (u *user) Login(ctx context.Context, login model.Login) (map[string]string, error) {
	userInstance := model.User{
		Username: login.Username,
		Password: login.Password,
	}

	checkAccess, id, err := u.LoginCheck(ctx, userInstance)
	if err != nil {
		return nil, fmt.Errorf("user.Login.LoginCheck: %v", err)
	}

	accessTokenString, err := checkAccess.GenerateAccessToken(id)
	refreshTokenString, err := checkAccess.GenerateRefreshToken(id)

	tokenPair := map[string]string{"access_token": accessTokenString, "refresh_token": refreshTokenString}

	return tokenPair, nil
}

func (u *user) LoginCheck(ctx context.Context, loginUser model.User) (Token, uint, error) {
	hashedPass, err := u.repo.GetHash(ctx, loginUser.Username)
	if err != nil {
		return nil, 0, fmt.Errorf("user.LoginCheck.GetHash: %v", err)
	}

	id, err := u.repo.GetID(ctx, loginUser.Username)
	if err != nil {
		return nil, 0, fmt.Errorf("user.LoginCheck.GetID: %v", err)
	}

	if err = VerifyPassword(ctx, loginUser.Password, hashedPass); err != nil {
		return nil, 0, fmt.Errorf("user.LoginCheck.VerifyPassword: %v", err)
	}

	return NewToken(os.Getenv("API_SECRET")), id, nil
}

func VerifyPassword(ctx context.Context, pass, hashPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(pass))
}
