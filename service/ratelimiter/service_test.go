package ratelimiter

import (
	"context"
	"testing"
	"time"

	"github.com/Mohamadreza-shad/ratelimiter/client"
	"github.com/Mohamadreza-shad/ratelimiter/config"

	"github.com/Mohamadreza-shad/ratelimiter/logger"
)

func BenchmarkRateLimit(b *testing.B) {
	userID := "1234"
	windowSize := time.Duration(3) * time.Second
	redisClient, err := client.NewRedisClient()
	if err != nil {
		panic(err)
	}
	err = redisClient.FlushAll(context.Background()).Err()
	if err != nil {
		panic(err)
	}
	logger, err := logger.New()
	if err != nil {
		panic(err)
	}
	rateLimiterService := New(redisClient, windowSize, logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rateLimiterService.RateLimit(userID, config.UserRateLimit())
	}
}
