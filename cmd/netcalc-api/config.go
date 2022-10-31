package main

import "github.com/spf13/viper"

type Config struct {
	HttpPort string `mapstructure:"HTTP_PORT"`
}

func LoadConfig() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil
	}

	config := &Config{}
	err = viper.Unmarshal(config)

	return config
}
