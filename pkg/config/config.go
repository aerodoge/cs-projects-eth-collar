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
	viper.SetDefault("monitor.account", "default")
	viper.SetDefault("prometheus.enabled", true)
	viper.SetDefault("prometheus.push_gateway.url", "http://localhost:9091")
	viper.SetDefault("prometheus.push_gateway.job_name", "deribit-monitor")
	viper.SetDefault("prometheus.push_gateway.instance", "default")
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
