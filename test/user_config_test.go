package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"web-metric/api"
	"web-metric/service/ratelimiter"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func Test_SetUserManualRateLimit_Successful(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)
	logger := getLogger()
	validator := validator.New()
	redisClient := getRedis()
	err := redisClient.FlushAll(ctx).Err()
	assert.Nil(err)
	windowSize := time.Duration(3) * time.Second
	rateLimitService := ratelimiter.New(redisClient, windowSize, logger)
	userHandler := api.NewUserHandler(rateLimitService, validator)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/api/v1/user/rate-limit", userHandler.SetUserManualRateLimit)

	server := httptest.NewServer(r)
	defer server.Close()

	params := ratelimiter.SetUserConfigParams{
		UserID:    "123",
		RateLimit: 5,
	}
	jsonData, err := json.Marshal(params)
	assert.Nil(err)

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/api/v1/user/rate-limit", server.URL),
		bytes.NewBuffer(jsonData),
	)
	assert.Nil(err)

	_, err = http.DefaultClient.Do(req)
	assert.Nil(err)

	key := "userID: 123"
	userManualRatelimiter, err := redisClient.HGet(ctx, key, ratelimiter.REDIS_RATE_LIMIT_FIELD).Int()
	assert.Nil(err)
	assert.Equal(userManualRatelimiter, int(5))
}
