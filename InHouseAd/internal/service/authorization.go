package service

import (
	"InHouseAd/internal/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Token interface {
	GenerateAccessToken(userId uint) (string, error)
	GenerateRefreshToken(userId uint) (string, error)
	ValidateToken(tokenString string) (*model.TokenClaims, error)
	RefreshToken(refreshToken string) (string, error)
	ValidateRefreshToken(c *gin.Context)
}

type token struct {
	apiSecret string
}

func NewToken(apiSecret string) Token {
	return &token{apiSecret: apiSecret}
}

var accessTokenMaxAge = time.Now().Add(time.Minute * 10).Unix()
var refreshTokenMaxAge = time.Now().Add(time.Minute * 20).Unix()

func (t *token) GenerateAccessToken(userId uint) (string, error) {
	claims := model.TokenClaims{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenMaxAge,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(t.apiSecret))
}

func (t *token) GenerateRefreshToken(userId uint) (string, error) {
	claims := model.TokenClaims{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenMaxAge,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(t.apiSecret))
}

func (t *token) ValidateToken(tokenString string) (*model.TokenClaims, error) {
	vToken, err := jwt.ParseWithClaims(tokenString, &model.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.apiSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("authorization.ValidateToken.ParseWithClaims: %v", err)
	}

	if claims, ok := vToken.Claims.(*model.TokenClaims); ok && vToken.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("authorization.ValidateToken: %v", err)
}

func (t *token) RefreshToken(refreshToken string) (string, error) {
	claims, err := t.ValidateToken(refreshToken)

	if err != nil {
		return "", err
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 0 {
		return "", fmt.Errorf("refresh token expired")
	}

	return t.GenerateAccessToken(claims.UserID)
}

func (t *token) ValidateRefreshToken(c *gin.Context) {
	rToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	newAccessToken, err := t.RefreshToken(rToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token unauthorized"})
		return
	}

	c.SetCookie("access_token", newAccessToken, int(accessTokenMaxAge), "/", "localhost", false, false)
}
