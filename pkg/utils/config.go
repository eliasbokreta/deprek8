package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

// Load configuration file
func LoadConfig() error {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Set configuration default values
	viper.SetDefault("deprek8.helm.repository.artifacthub.apiBaseUrl", "https://artifacthub.io/api/v1")
	viper.SetDefault("deprek8.helm.resourcesFile", "./config/api.yaml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("could not retrieve config file: %w", err)
	}

	return nil
}
