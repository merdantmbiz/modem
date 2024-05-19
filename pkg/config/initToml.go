package config

import (
	"github.com/spf13/viper"
)

var TomlConf Config

func InitTomlConf(configFileName, configPath string) error {
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(configPath)
	viper.SetConfigType("toml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		
		return err
	}
	
	err = viper.Unmarshal(&TomlConf)
	if err != nil {
		return err
	}
	return nil
}