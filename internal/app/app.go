package app

import (
	"InHouseAd/docs"
	"InHouseAd/internal/config"
	"InHouseAd/internal/fetch"
	"InHouseAd/internal/handler"
	"InHouseAd/internal/middleware"
	"InHouseAd/internal/pkg"
	"InHouseAd/internal/repository"
	"InHouseAd/internal/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"sync"
	"time"
)

func Run(cfg *config.Config) {
	ctx := context.Background()

	logger := pkg.NewLogger(cfg)
	newLog, err := logger.InitLogger(ctx)
	if err != nil {
		log.Panic(err)
	}

	pgSession := pkg.NewPgSession(&cfg.Postgres)
	pool, err := pgSession.InitPgSession(ctx)
	if err != nil {
		log.Panic(err)
	}

	userRepo := repository.NewUser(pool)
	goodRepo := repository.GoodRepo(pool)
	categoryRepo := repository.CategoryRepo(pool)

	userService := service.NewUser(userRepo)
	goodService := service.GoodService(goodRepo)
	categoryService := service.CategoryService(categoryRepo, goodRepo)

	userHandler := handler.NewUserHandler(newLog, userService)
	categoryHandler := handler.CategoryHandler(categoryService, newLog)
	goodHandler := handler.GoodHandler(goodService, newLog)

	auth := service.NewToken(cfg.ApiSecret)
	mware := middleware.NewMiddleware(newLog, auth)

	var timer time.Timer
	interval := time.Hour * 1
	goodsFetch := fetch.NewFetch(goodRepo, categoryRepo, &timer)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		goodsFetch.Init(ctx, interval)
	}()
	wg.Wait()

	router := gin.New()
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.ServerAddress, cfg.ServerPort)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	goodsActions := router.Group("/good")
	{
		goodsActions.GET("/list/:id", goodHandler.GetByCategory)
		goodsActions.POST("/create", mware.JwtAuthMiddleware, auth.ValidateRefreshToken, goodHandler.Create)
		goodsActions.DELETE("/remove/:id", mware.JwtAuthMiddleware, auth.ValidateRefreshToken, goodHandler.Delete)
		goodsActions.PATCH("/update/:id", mware.JwtAuthMiddleware, auth.ValidateRefreshToken, goodHandler.Update)
	}

	categoryActions := router.Group("/category")
	{
		categoryActions.GET("/list", categoryHandler.Get)
		categoryActions.POST("/create", mware.JwtAuthMiddleware, auth.ValidateRefreshToken, categoryHandler.Create)
		categoryActions.DELETE("/remove/:id", mware.JwtAuthMiddleware, auth.ValidateRefreshToken, categoryHandler.Delete)
		categoryActions.PATCH("/update/:id", mware.JwtAuthMiddleware, auth.ValidateRefreshToken, categoryHandler.Update)
	}

	if err := router.Run(docs.SwaggerInfo.Host); err != nil {
		log.Panic(err)
	}
}
