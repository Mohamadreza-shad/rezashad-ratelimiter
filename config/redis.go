package config

type RedisMasterNameConfigs struct {
	Name string
}

type RedisConfigs struct {
	URI    string
	Master RedisMasterNameConfigs
}

func RedisURI() string {
	if cfg.Redis.URI == "" {
		return "redis://:123456@localhost:6379"
	}
	return cfg.Redis.URI
}

func RedisMasterName() string {
	return cfg.Redis.Master.Name
}
