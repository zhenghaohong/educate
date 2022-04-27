package config

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"fmt"

	"github.com/spf13/viper"
)

var (
	cfgMqtt        *viper.Viper
	cfgWs          *viper.Viper
	cfgTopic       *viper.Viper
	cfgDatabase    *viper.Viper
	cfgRedis 	   *viper.Viper
	cfgDeviceConfig *viper.Viper
	
)

func init() {
	Setup("config/config.yml")
}

//载入配置文件
func Setup(path string) {
	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}

	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		log.Fatal(fmt.Sprintf("Parse config file fail: %s", err.Error()))
	}

	cfgMqtt = viper.Sub("settings.mqtt")
	if cfgMqtt == nil {
		panic("No found settings.mqtt in the configuration")
	}
	MQTTConfig = InitMQTT(cfgMqtt)

	cfgWs = viper.Sub("settings.ws")
	if cfgWs == nil {
		panic("No found settings.ws in the configuration")
	}
	WsConfig = InitWs(cfgWs)

	// cfgHttpApi = viper.Sub("settings.httpapi")
	// if cfgHttpApi == nil {
	// 	panic("No found settings.httpapi in the configuration")
	// }

	cfgTopic = viper.Sub("settings.topic")
	if cfgTopic == nil {
		panic("No found settings.topic in the configuration")
	}
	TopicConfig = InitTopic(cfgTopic)

	cfgDatabase = viper.Sub("settings.database")
	if cfgDatabase == nil {
		panic("No found settings.database in the configuration")
	}
	DatabaseConfig = InitDatabase(cfgDatabase)


	cfgRedis = viper.Sub("settings.redis")
	if cfgRedis == nil {
		panic("No found settings.redis in the configuration")
	}
	RedisConfig = InitRedis(cfgRedis)

	cfgDeviceConfig = viper.Sub("settings.deviceControlConfig")
	if cfgDeviceConfig == nil {
		panic("No found settings.deviceControlConfig in the configuration")
	}
	DeviceConfigConfig = InitDeviceConfig(cfgDeviceConfig)
}
