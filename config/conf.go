package config

import (
	"github.com/spf13/viper"
	"log"
)

type RestConfig struct {
	Host string
	Port int
}

type AppConfig struct {
	Env             string
	GracefulTimeout int
	Rest            RestConfig
}

var appConfig *AppConfig

func LoadConfig() error {
	appConfig = new(AppConfig)
	if err := viper.UnmarshalKey("config", appConfig); err != nil {
		log.Println("Server could not Load config:", err)
		return err
	}
	return nil
}

func GetConfig() *AppConfig {

	return appConfig
}
