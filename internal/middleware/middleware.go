package middleware

import (
	"InHouseAd/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
)

type Middleware interface {
	JwtAuthMiddleware(c *gin.Context)
}

type middleware struct {
	logger zerolog.Logger
	tkn    service.Token
}

func NewMiddleware(logger zerolog.Logger, tkn service.Token) Middleware {
	return &middleware{tkn: tkn, logger: logger}
}

func (m *middleware) JwtAuthMiddleware(c *gin.Context) {
	tokenString, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		m.logger.Error().Msgf("no token found: %v", fmt.Errorf("middleware.JwtAuthMiddleware.Cookie: %v", err))
		c.Abort()
		return
	}

	_, err = m.tkn.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		m.logger.Error().Msgf("no token found: %v", fmt.Errorf("middleware.JwtAuthMiddleware.ValidateToken: %v", err))
		c.Abort()
		return
	}

	c.Next()
}
