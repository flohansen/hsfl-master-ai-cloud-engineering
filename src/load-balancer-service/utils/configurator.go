package utils

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Configuration struct {
	Image    string `yaml:"image"`    // Image which should be run.
	Replicas int    `yaml:"replicas"` // Replicas number to run and load balance.
	Network  string `yaml:"network"`  // Network where replicas belong to.
}

const (
	DefaultConfigPath string = "./config/config.yaml"
)

func GetConfig(path string) *Configuration {
	viper.SetDefault("image", "nginxdemos/hello")
	viper.SetDefault("replicas", 5)
	viper.SetDefault("network", "bridge")

	if len(path) != 0 {
		viper.SetConfigFile(path)
	} else {
		if _, err := os.Stat(DefaultConfigPath); errors.Is(err, os.ErrNotExist) {
			viper.SetConfigType("yaml")
			viper.WriteConfigAs(DefaultConfigPath)
		}
		viper.SetConfigFile(DefaultConfigPath)
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	configuration := &Configuration{}
	if err := viper.Unmarshal(configuration); err != nil {
		fmt.Printf("unable to decode configuration, %v", err)
	}

	return configuration
}
