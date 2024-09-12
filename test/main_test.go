package test

import (
	"os"
	"testing"
	"web-metric/client"
	"web-metric/config"
	"web-metric/logger"

	"github.com/redis/go-redis/v9"
)

var loggerService *logger.Logger
var redisClient redis.UniversalClient

func TestMain(m *testing.M) {
	config.Load()
	if config.Env() != config.EnvTest {
		config.SetTestEnvVariable()
	}
	exitCode := m.Run()
	os.Exit(exitCode)
}

func getLogger() *logger.Logger {
	if loggerService != nil {
		return loggerService
	}
	var err error
	loggerService, err = logger.New()
	if err != nil {
		panic(err)
	}
	return loggerService
}

func getRedis() redis.UniversalClient {
	if redisClient != nil {
		return redisClient
	}
	redisClient, _ = client.NewRedisClient()
	return redisClient
}
