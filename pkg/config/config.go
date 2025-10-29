package config

import (
	"cs-projects-eth-collar/internal/types"
	"github.com/spf13/viper"
)

func LoadConfig(configPath string) (*types.Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	viper.SetDefault("deribit.base_url", "https://www.deribit.com/api/v2")
	viper.SetDefault("deribit.test_net", false)
	viper.SetDefault("monitor.interval_seconds", 30)
	viper.SetDefault("monitor.mm_threshold", 0.5)
	viper.SetDefault("monitor.mm_target", 0.3)
	viper.SetDefault("monitor.eth_equity_threshold", -700000.0)
	viper.SetDefault("monitor.eth_equity_target", 200.0)
	viper.SetDefault("alerts.enabled", true)
	viper.SetDefault("alerts.methods", []string{"log"})
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.file", "monitor.log")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config types.Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
