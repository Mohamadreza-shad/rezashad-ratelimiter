package ratelimiter

import (
	"context"
	"errors"
	"time"
	"web-metric/config"
	"web-metric/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap/zapcore"
)

const (
	REDIS_RATE_LIMIT_FIELD = "rate_limit"
)

type Service struct {
	redisClient redis.UniversalClient
	windowSize  time.Duration
	logger      *logger.Logger
}

//The RateLimit method implements a rate-limiting mechanism using Redis. 
//It limits the number of requests a user can make within a specified time window. 
//The method ensures that requests beyond the defined limit are blocked, 
//enforcing fair usage of the system's resources.
func (s *Service) RateLimit(userID string, limit int) bool {
	ctx := context.Background()

	// Define the Redis key for this user's rate limit
	key := "user:" + userID + ":limit"

	// Fetch the user-specific limit from Redis (optional)
	userLimit, err := s.redisClient.HGet(ctx, "userID:"+userID, REDIS_RATE_LIMIT_FIELD).Int()
	if err == redis.Nil {
		userLimit = 1 // Default limit if not set for the user
	} else if err != nil {
		s.logger.Log(zapcore.ErrorLevel, "failed to get user limit. err: "+err.Error())
		return false
	}

	// Total limit (default + user-specific limit)
	totalLimit := limit + userLimit

	// Increment the count for this user in Redis atomically
	count, err := s.redisClient.Incr(ctx, key).Result()
	if err != nil {
		s.logger.Log(zapcore.ErrorLevel, "failed to increment rate limit counter. err: "+err.Error())
		return false
	}

	// If it's the first request, set an expiration (TTL) on the key
	if count == 1 {
		err := s.redisClient.Expire(ctx, key, s.windowSize).Err()
		if err != nil {
			s.logger.Log(zapcore.ErrorLevel, "failed to set expiration on rate limit key. err: "+err.Error())
			return false
		}
	}

	// Deny the request if the count exceeds the limit
	if int(count) > totalLimit {
		return false
	}

	// Allow the request
	return true
}

// The SetUserConfig method allows the system to configure and set custom
// rate limits for individual users in Redis. 
// This method is typically used to update or create
// user-specific rate limit configurations that override the default rate limit.
func (s *Service) SetUserConfig(ctx context.Context, params SetUserConfigParams) error {
	key := "userID:" + params.UserID
	if params.RateLimit == 0 {
		params.RateLimit = config.UserRateDefault()
	}
	err := s.redisClient.HSet(ctx, key, REDIS_RATE_LIMIT_FIELD, params.RateLimit).Err()
	if err != nil {
		s.logger.Log(zapcore.ErrorLevel, "failed to set user config")
		return errors.New("failed to set user config")
	}
	return nil
}


func (s *Service) AllProducts(ctx context.Context) []Product {
	return []Product{
		{
			ID:   1,
			Name: "Product1",
		},
		{
			ID:   2,
			Name: "Product2",
		},
		{
			ID:   3,
			Name: "Product3",
		},
	}
}

func New(
	redisClient redis.UniversalClient,
	windowSize time.Duration,
	logger *logger.Logger,
) *Service {
	return &Service{
		redisClient: redisClient,
		windowSize:  windowSize,
		logger:      logger,
	}
}
