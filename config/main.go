package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	Minio MinioConfig
}

// AllConfig variable of type AppConfig
var AllConfig AppConfig

// GetConfig Collects all configs
func GetConfig() AppConfig {
	_ = godotenv.Load()

	AllConfig = AppConfig{}

	err := envconfig.Process("APP_PORT", &AllConfig)
	if err != nil {
		panic(err)
	}

	return AllConfig
}

// GetConfigByName Collects all configs
func GetConfigByName(key string) string {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	return os.Getenv(key)
}
