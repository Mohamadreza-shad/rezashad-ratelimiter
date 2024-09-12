package config

import (
	"flag"

	"github.com/pkg/errors"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/source/env"
)

const (
	EnvProd = "prod"
	EnvTest = "test"
	EnvDev  = "dev"
)

var cfg *Config = &Config{}

type Config struct {
	Environment string //related to sentry
	Env         string
	Redis       RedisConfigs
	Hostname    string
	UserRate    UserRate
	Window      Window
	Server      Server
}

func Env() string {
	return cfg.Env
}

func isTestEnv() bool {
	return flag.Lookup("test.v") != nil
}

func SetTestEnvVariable() {
	cfg.Env = EnvTest
}

func Load() error {
	config, err := config.NewConfig(config.WithSource(env.NewSource()))
	if err != nil {
		return errors.Wrap(err, "config.New")
	}
	if err := config.Load(); err != nil {
		return errors.Wrap(err, "config.Load")
	}
	if err := config.Scan(cfg); err != nil {
		return errors.Wrap(err, "config.Scan")
	}
	if isTestEnv() {
		SetTestEnvVariable()
	}
	return nil
}
