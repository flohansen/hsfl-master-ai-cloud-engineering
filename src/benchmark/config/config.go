package config

import (
	"encoding/json"
	"os"
)

type LoadTestConfig struct {
	Users    int      `json:"users"`
	RampUp   float64  `json:"rampup"`
	Duration float64  `json:"duration"`
	Targets  []string `json:"targets"`
}

func LoadConfig(path string) (*LoadTestConfig, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var config LoadTestConfig

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
