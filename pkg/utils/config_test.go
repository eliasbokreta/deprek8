package utils

import (
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfig(t *testing.T) {
	viper.AddConfigPath("../../config")

	if err := LoadConfig(); err != nil {
		t.Fatalf("Error trying to load the config: %v", err)
	}

	if viper.GetString("deprek8.helm.repository.artifacthub.apiBaseUrl") != "https://artifacthub.io/api/v1" {
		t.Fatal("'apiBaseUrl' config should be equal to 'https://artifacthub.io/api/v1'")
	}

	if viper.GetString("deprek8.helm.resourcesFile") != "./config/api.yaml" {
		t.Fatal("'resourcesFile' config should be equal to './config/api.yaml'")
	}
}
