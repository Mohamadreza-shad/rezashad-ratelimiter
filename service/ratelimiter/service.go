package ratelimiter

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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

type SetUserConfigParams struct {
	UserID    string `json:"userID"`
	RateLimit int    `json:"rateLimit"`
}

// Please go to README.md file for furthur description
func (s *Service) RateLimit(userID string, limit int) bool {
	ctx := context.Background()
	currentTime := time.Now().UnixMilli()
	windowStart := currentTime - s.windowSize.Milliseconds()

	// Define the Redis key for this user
	key := "user:" + userID + "limit"

	// Remove old requests that are outside the sliding window
	s.redisClient.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))

	// Count the requests within the current window
	count, err := s.redisClient.ZCount(ctx, key, strconv.FormatInt(windowStart, 10), strconv.FormatInt(currentTime, 10)).Result()
	if err != nil {
		return false
	}

	//check if user has more limit in his/her config
	userLimit, err := s.redisClient.HGet(ctx, "userID:"+userID, REDIS_RATE_LIMIT_FIELD).Int()
	if err == redis.Nil {
		s.logger.Log(zapcore.ErrorLevel, "manual rate-limiter has not been set for user")
		err := s.redisClient.HSet(ctx, "userID:"+userID, REDIS_RATE_LIMIT_FIELD, config.UserRateDefault()).Err()
		if err != nil {
			s.logger.Log(zapcore.ErrorLevel, "failed to set user config")
		}
	}
	fmt.Println("count: ", count, "limit: ", int(limit+userLimit))
	// Deny request as the user has exceeded the limit
	if int(count) >= int(limit+userLimit) {
		return false
	}

	// Otherwise, allow the request and add the current timestamp to Redis
	s.redisClient.ZAdd(ctx, key, redis.Z{Score: float64(currentTime), Member: currentTime})

	// Optionally set a TTL for the sorted set to auto-expire
	s.redisClient.Expire(ctx, key, s.windowSize)

	return true
}

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

type Product struct {
	ID   int
	Name string
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
