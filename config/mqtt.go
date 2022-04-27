package config
import "github.com/spf13/viper"

type Mqtt struct {
	Username   string
	Password   string
	ClientId   string
	IP         string
	Port       int
}

func InitMQTT(cfg *viper.Viper) *Mqtt {
	return &Mqtt{
		Username:   cfg.GetString("username"),
		Password:   cfg.GetString("password"),
		ClientId:   cfg.GetString("clientId"),
		IP:         cfg.GetString("ip"),
		Port:       cfg.GetInt("port"),
	}
}

var MQTTConfig = new(Mqtt)

