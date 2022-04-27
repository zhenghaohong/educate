package config

import "github.com/spf13/viper"

type Redis struct {
	Addr string
	Password string
}

func InitRedis(cfg *viper.Viper) *Redis {

	db := &Redis{
		Addr: cfg.GetString("addr"),
		Password: cfg.GetString("password"),
	}
	return db
}

var RedisConfig = new(Redis)