package config

type Users struct {
	Rate Rate
}

type Rate struct {
	Limit   int
	Default int
}

func UsersRateLimit() int {
	if cfg.Users.Rate.Limit == 0 && Env() == EnvTest{
		return 198
	}
	return cfg.Users.Rate.Limit
}

func UsersRateDefault() int {
	if cfg.Users.Rate.Default == 0 {
		return 1
	}
	return cfg.Users.Rate.Default
}
