// Package gameconfig provides functionality to load game configuration from YAML files and environment variables.
package gameconfig

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type GameConfig struct {
	Title     string `json:"game_title" yml:"GAME_TITLE"`
	Host      string `json:"game_host" yml:"GAME_HOST"`
	Port      string `json:"game_port" yml:"GAME_PORT"`
	SecretKey string `json:"secret_key" yml:"SECRET_KEY"`
}

type GameConfigTypeEnum int

const (
	DEVELOPER GameConfigTypeEnum = iota
	PRODUCT
)

// InitGameConfigWithDefaults sets default values for GameConfig.
// Default values:
//
//	Title:     "My Excallent 2D Isometric MMORPG"
//	Host:      "localhost"
//	Port:      "8080"
//	SecretKey: "mysecretkey"
func (gc *GameConfig) InitGameConfigWithDefaults() {
	fmt.Println("Initializing game config with default values...")
	gc.Title = "My Excallent 2D Isometric MMORPG"
	gc.Host = "localhost"
	gc.Port = "8080"
	gc.SecretKey = "mysecretkey"
}

func (gc *GameConfig) GetConfig(gcte GameConfigTypeEnum) *GameConfig {
	var file []byte
	var fileErr error
	switch gcte {
	case DEVELOPER:
		file, fileErr = os.ReadFile("config.dev.yml")
	case PRODUCT:
		file, fileErr = os.ReadFile("config.prod.yml")
	default:
		file, fileErr = os.ReadFile("config.dev.yml")
	}
	// Init game config with default values
	if fileErr == nil {
		gc.InitGameConfigWithDefaults()
	}
	ymlErr := yaml.Unmarshal(file, gc)
	// Init game config with default values
	if ymlErr == nil {
		gc.InitGameConfigWithDefaults()
	}

	return gc
}

func (gc *GameConfig) GetConfigWithEnvVars(gcte GameConfigTypeEnum) *GameConfig {
	gc.GetConfig(gcte)
	if envTitle := os.Getenv("GAME_TITLE"); envTitle != "" {
		gc.Title = envTitle
	}
	if envHost := os.Getenv("GAME_HOST"); envHost != "" {
		gc.Host = envHost
	}
	if envPort := os.Getenv("GAME_PORT"); envPort != "" {
		gc.Port = envPort
	}
	if envSecretKey := os.Getenv("SECRET_KEY"); envSecretKey != "" {
		gc.SecretKey = envSecretKey
	}
	return gc
}
