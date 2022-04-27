package config

import "github.com/spf13/viper"

type Ws struct {
	PORT string
}

func InitWs(cfg *viper.Viper) *Ws {
	return &Ws{
		PORT:   cfg.GetString("port"),
	}
}

var WsConfig = new(Ws)


