package config

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"os"
	"path/filepath"
)

func GetConfig() (Config, error) {

	env := GetEnvironment()
	var cfg Config
	fileConfig := "/etc/secrets/config.json"
	if env != "prod" && env != "stage" {
		dir, err := filepath.Abs("./")
		if err != nil {
			return Config{}, err
		}
		fileConfig = filepath.Join(dir, "/app/shared/utils/config/config.json")
	}

	byteConfig, err := os.ReadFile(fileConfig)
	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(byteConfig, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing config")
	}

	validate := validator.New()
	err = validate.Struct(cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func GetEnvironment() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "LOCAL"
	}
	return env
}
