package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type RestConfig struct {
	Host string
	Port int
}

type UploadConfig struct {
	BaseUrl  string
	S3Bucket string
	Folder   string
}

type ApiConfig struct {
	AuthApi string
}

type AppConfig struct {
	Env             string
	GracefulTimeout int
	Rest            RestConfig
	Upload          UploadConfig
	APIs            ApiConfig
}

var appConfig *AppConfig

func LoadConfig() error {
	appConfig = new(AppConfig)
	if err := viper.UnmarshalKey("config", appConfig); err != nil {
		log.Println("Server could not Load config:", err)
		return err
	}
	err := os.Setenv("AUTH_HOST", appConfig.APIs.AuthApi)
	if err != nil {
		return err
	}
	return nil
}

func GetConfig() *AppConfig {

	return appConfig
}
