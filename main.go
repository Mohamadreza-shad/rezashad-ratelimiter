package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Mohamadreza-shad/ratelimiter/api"
	"github.com/Mohamadreza-shad/ratelimiter/api/router"
	"github.com/Mohamadreza-shad/ratelimiter/client"
	"github.com/Mohamadreza-shad/ratelimiter/config"
	"github.com/Mohamadreza-shad/ratelimiter/logger"
	"github.com/Mohamadreza-shad/ratelimiter/service/ratelimiter"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func main() {
	//initializing logger
	logger, err := logger.New()
	if err != nil {
		log.Fatal("failed to initialize logger", err)
	}
	
	//loading config
	err = config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}
	defer logger.Sync()
	fmt.Println(config.UsersRateLimit())
	//initiating redis client
	redisClient, err := client.NewRedisClient()
	if err != nil {
		logger.Fatal("failed to initiate redis client. " + err.Error())
	}

	//initiating rate-limiter service
	windowSize := 10 * time.Duration(config.WindowSize()) * time.Second
	rateLimiterService := ratelimiter.New(redisClient, windowSize, logger)
	validator := validator.New()
	userHandler := api.NewUserHandler(rateLimiterService, validator)

	//initiating router object
	router := router.New(rateLimiterService, userHandler, logger)
	httpServer := &http.Server{
		Addr:    config.ServerHttpAddress(),
		Handler: router.Handler,
	}
	logger.Info(
		"starting HTTP server on %s",
		zap.String("HTTP server address: ",
			config.ServerHttpAddress()),
	)

	//preparing the server
	err = httpServer.ListenAndServe()
	if err != nil {
		logger.Fatal(err.Error())
	}
}
