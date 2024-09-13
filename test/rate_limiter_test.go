package test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Mohamadreza-shad/ratelimiter/api"
	"github.com/Mohamadreza-shad/ratelimiter/api/middleware"
	"github.com/Mohamadreza-shad/ratelimiter/service/ratelimiter"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_1LimitForManualAnd3ForLimit_OneReqShouldGet429(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)
	logger := getLogger()
	validator := validator.New()
	redisClient := getRedis()
	err := redisClient.FlushAll(ctx).Err()
	assert.Nil(err)

	windowSize := time.Duration(3) * time.Second
	currentTime := time.Now().UnixMilli()
	windowStart := currentTime - windowSize.Milliseconds()
	rateLimitService := ratelimiter.New(redisClient, windowSize, logger)
	userHandler := api.NewUserHandler(rateLimitService, validator)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(middleware.RateLimitMiddleware(rateLimitService))
	r.GET("/api/v1/products", userHandler.AllProducts)

	server := httptest.NewServer(r)
	defer server.Close()

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/api/v1/products?userID=123", server.URL),
		nil,
	)
	assert.Nil(err)

	concurrentRequests :=200 
	ch := make(chan int, concurrentRequests)
	var wg sync.WaitGroup
	wg.Add(concurrentRequests)
	for i := 0; i < concurrentRequests; i++ {
		go func() {
			res, err := http.DefaultClient.Do(req)
			assert.Nil(err)
			ch <- int(res.StatusCode)
			defer wg.Done()
		}()
	}
	wg.Wait()
	close(ch)
	blockedReq := []int{}
	for c := range ch {
		if c == int(http.StatusTooManyRequests) {
			blockedReq = append(blockedReq, c)
		}
	}
	assert.Equal(len(blockedReq), int(1))
	assert.Nil(err)

	userID := "123"
	key := "user:" + userID + "limit"
	count, err := redisClient.ZCount(ctx, key, strconv.FormatInt(windowStart, 10), strconv.FormatInt(currentTime, 10)).Result()
	assert.Nil(err)
	fmt.Println(count)
}

func Test_5ConcurrentReques_ThenWaitForWindowToClose_ThenCallAgain_WeShouldGet200InsteadOf429(t *testing.T) {
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
	r.Use(middleware.RateLimitMiddleware(rateLimitService))
	r.GET("/api/v1/products", userHandler.AllProducts)

	server := httptest.NewServer(r)
	defer server.Close()

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/api/v1/products?userID=123", server.URL),
		nil,
	)
	assert.Nil(err)

	concurrentRequests := 200
	ch := make(chan int, concurrentRequests)
	var wg sync.WaitGroup
	wg.Add(concurrentRequests)
	for i := 0; i < concurrentRequests; i++ {
		go func() {
			res, err := http.DefaultClient.Do(req)
			assert.Nil(err)
			ch <- int(res.StatusCode)
			defer wg.Done()
		}()
	}
	wg.Wait()
	close(ch)
	blockedReq := []int{}
	for c := range ch {
		if c == int(http.StatusTooManyRequests) {
			blockedReq = append(blockedReq, c)
		}
	}
	assert.Equal(len(blockedReq), int(1))
	assert.Nil(err)

	time.Sleep(4 * time.Second) //wait for window time to pass and try again. next request shoul get 200 instead of 429
	res, err := http.DefaultClient.Do(req)
	assert.Nil(err)
	assert.Equal(int(res.StatusCode), int(200))
}
