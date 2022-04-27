package config

import "github.com/spf13/viper"

type DeviceConfig struct {
	SupineConfigFarthestDistance int
	SupineConfigShortestDistance int
	SupineConfigLastTime int
	JumpConfigUnitTime int

	JumpConfigUnitValue int
	PullUpFarthestDistance int
	PullUpShortestDistance int
	PullUpLastTime int
	SupineFarthestDistance int 
	SupineShortestDistance int
	SupineLastTime int
}

func InitDeviceConfig(cfg *viper.Viper) *DeviceConfig {
	db := &DeviceConfig{
		SupineConfigFarthestDistance: cfg.GetInt("SupineConfigFarthestDistance"),
		SupineConfigShortestDistance: cfg.GetInt("SupineConfigShortestDistance"),
		SupineConfigLastTime: cfg.GetInt("SupineConfigLastTime"),
		JumpConfigUnitTime: cfg.GetInt("JumpConfigUnitTime"),
		JumpConfigUnitValue: cfg.GetInt("JumpConfigUnitValue"),
		PullUpFarthestDistance: cfg.GetInt("PullUpFarthestDistance"),
		PullUpShortestDistance: cfg.GetInt("PullUpShortestDistance"),
		PullUpLastTime: cfg.GetInt("PullUpLastTime"),
		SupineFarthestDistance: cfg.GetInt("SupineFarthestDistance"),
		SupineShortestDistance: cfg.GetInt("SupineShortestDistance"),
		SupineLastTime: cfg.GetInt("SupineLastTime"),
	}
	return db
}

var DeviceConfigConfig = new(DeviceConfig)
