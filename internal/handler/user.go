package handler

import (
	"InHouseAd/internal/model"
	"InHouseAd/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type User interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

// user interface struct
type user struct {
	logger      zerolog.Logger
	userService service.User
}

// NewUserHandler constructor-method
func NewUserHandler(logger zerolog.Logger, userService service.User) User {
	return &user{logger: logger, userService: userService}
}

var accessTokenMaxAge = time.Now().Add(time.Minute * 1).Unix()
var refreshTokenMaxAge = time.Now().Add(time.Hour * 3).Unix()

// Register Handle function for registration request
// @Summary Registration-endpoint
// @Description User-registration
// @Produce json
// @Param regInput body model.Register true "Registration data"
// @Success 200 {object} model.User
// @Router /register [post]
func (a *user) Register(c *gin.Context) {
	var regUser model.Register
	if err := c.ShouldBindJSON(&regUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		a.logger.Error().Msgf("Bad Request : %v", err)
		return
	}

	err := a.userService.Register(c, regUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not register"})
		a.logger.Error().Msgf("error registering new user : %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
	a.logger.Error().Msg("registration successful")
}

// Login Handle function for login request
// @Summary Login-endpoint
// @Description User-login
// @Produce json
// @Param loginInput body model.Login true "Login data"
// @Success 200 {object} model.User
// @Router /login [post]
func (a *user) Login(c *gin.Context) {
	var loginUser model.Login
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		a.logger.Error().Msgf("Bad Request : %v", err)
		return
	}

	tokenPair, err := a.userService.Login(c, loginUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		a.logger.Error().Msgf("error logging in: %v", err)
		return
	}

	c.SetCookie("access_token", tokenPair["access_token"], int(accessTokenMaxAge), "/", "localhost", false, false)
	c.SetCookie("refresh_token", tokenPair["refresh_token"], int(refreshTokenMaxAge), "/", "localhost", false, false)

	c.JSON(http.StatusOK, gin.H{"token": tokenPair["access_token"]})
	a.logger.Error().Msg("login successful")
}
