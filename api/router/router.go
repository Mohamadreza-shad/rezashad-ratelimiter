package router

import (
	"errors"
	"fmt"
	"net/http"
	"web-metric/api"
	"web-metric/api/middleware"
	"web-metric/logger"
	"web-metric/service/ratelimiter"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	Service = "gateway-service"
)

type Router struct {
	Handler *gin.Engine
}

func New(
	rateLimitService *ratelimiter.Service,
	userHandler *api.UserHandler,
	logger *logger.Logger,
) *Router {
	gin.SetMode(gin.ReleaseMode)
	router := &Router{}
	r := gin.New()

	r.Use(cors.Default())
	r.Use(globalRecover(logger))
	r.NoRoute(func(c *gin.Context) {
		c.JSON(
			http.StatusNotFound,
			api.ResponseFailure{
				Success: false,
				Error: api.ErrorCode{
					Code:    http.StatusNotFound,
					Message: "URL not found",
				},
			})
	})

	v1 := r.Group("/api/v1")
	v1.POST("/user/rate-limit", userHandler.SetUserManualRateLimit)

	securedV1 := r.Group("/api/v1")
	//implementing rate-limiter as a middleware
	securedV1.Use(middleware.RateLimitMiddleware(rateLimitService))
	securedV1.GET("/products", userHandler.AllProducts)

	router.Handler = r
	return router
}

func globalRecover(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(c *gin.Context) {
			if rec := recover(); rec != nil {
				err := errors.New("error 500")
				if err != nil {
					logger.Log(zapcore.ErrorLevel, fmt.Sprintf("error  500 in global recover %v", rec),
						zap.String("service", "httpServer"),
						zap.String("method", "globalRecover"),
					)
				}
				api.MakeErrorResponseWithCode(c.Writer, http.StatusInternalServerError, "error 500")
			}
		}(c)
		c.Next()
	}
}
