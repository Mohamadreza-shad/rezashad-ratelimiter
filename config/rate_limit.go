package config

type UserRate struct {
	Limit   int
	Default int
}

func UserRateLimit() int {
	if cfg.UserRate.Limit == 0 {
		return 3
	}
	return cfg.UserRate.Limit
}

func UserRateDefault() int {
	if cfg.UserRate.Default == 0 {
		return 1
	}
	return cfg.UserRate.Default
}
