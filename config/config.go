package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Database struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		DbName   string `mapstructure:"db_name"`
	} `mapstructure:"database"`
	Server struct {
		Port string `mapstructure:"port"`
	}
}

var config Config

func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	// Get the absolute path to the config directory
	configPath, err := filepath.Abs("./config")
	if err != nil {
		return nil, fmt.Errorf("error getting config path: %v", err)
	}

	// Set defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.username", "postgres")
	viper.SetDefault("database.db_name", "postgres")
	viper.SetDefault("server.port", "8080")

	viper.SetConfigType("yaml")
	viper.SetConfigName(env)
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")

	// Try to read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %v", err)
	}
	fmt.Println("The app is running on " + env)

	return &config, nil
}
