package model

import "github.com/dgrijalva/jwt-go"

type TokenClaims struct {
	UserID uint
	jwt.StandardClaims
}
