package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Target struct {
	URL string
}

type Configuration struct {
	Users     int       `json:"users"`     // Users to simulate or workers to run concurrently.
	Requests  int       `json:"requests"`  // Requests amount to send per user.
	RateLimit int       `json:"rateLimit"` // RateLimit for requests to send per second.
	Duration  int       `json:"duration"`  // Duration maximum of the test in seconds.
	RampUp    int       `json:"rampUp"`    // RampUp time for the ramping in seconds.
	Targets   []*Target `json:"targets"`   // Targets to address randomly
}

func GetConfig(path string) (*Configuration, error) {
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
		return nil, err
	}

	return configuration, nil
}
