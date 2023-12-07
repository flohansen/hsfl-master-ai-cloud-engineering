package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Target struct {
	URL string
}

type Configuration struct {
	Users     int       `yaml:"users"`      // Users to simulate or workers to run concurrently.
	Requests  int       `yaml:"requests"`   // Requests amount to send per user.
	RateLimit int       `yaml:"rate-limit"` // RateLimit for requests to send per second.
	Duration  int       `yaml:"duration"`   // Duration maximum of the test in seconds.
	RampUp    int       `yaml:"ramp-up"`    // RampUp time for the ramping in seconds.
	Targets   []*Target `yaml:"targets"`    // Targets to address randomly
}

func GetConfig(path string) *Configuration {
	viper.SetConfigFile(path)
	viper.SetDefault("users", 10)
	viper.SetDefault("requests", 1000)
	viper.SetDefault("rateLimit", 50)
	viper.SetDefault("duration", 100)
	viper.SetDefault("rampUp", 10)
	viper.SetDefault("targets", []*Target{
		{URL: "https://google.de:443"},
	})

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
