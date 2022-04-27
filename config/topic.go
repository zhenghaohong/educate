package config

import "github.com/spf13/viper"

type Topic struct {
	DeskDataTopic string
	DeskConfigTopic string
	SleepDataSubscribeTopic string
	SleepDataPublicTopic string
}

func InitTopic(cfg *viper.Viper) *Topic {
	return &Topic{
		DeskDataTopic:   cfg.GetString("deskDataTopic"),
		DeskConfigTopic:   cfg.GetString("deskConfigTopic"),
		SleepDataSubscribeTopic:   cfg.GetString("sleepDataSubscribeTopic"),
		SleepDataPublicTopic:   cfg.GetString("sleepDataPublicTopic"),
	}
}

var TopicConfig = new(Topic)

