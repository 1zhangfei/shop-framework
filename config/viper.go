package config

import "github.com/spf13/viper"

func ViperInit(address string) error {
	viper.SetConfigFile(address)
	return viper.ReadInConfig()
}
