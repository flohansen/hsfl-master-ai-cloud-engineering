package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type LoadTestConfig struct {
	Users         int           `json:"users"`
	StartSleep    time.Duration `json:"-"`
	StartSleepStr string        `json:"startSleep"`
	Specs         []Spec        `json:"specs"`
	Targets       []string      `json:"targets"`
}

type Spec struct {
	TargetSleep    time.Duration `json:"-"`
	TargetSleepStr string        `json:"targetSleep"`
	Duration       time.Duration `json:"-"`
	DurationStr    string        `json:"duration"`
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

	if len(config.Specs) < 1 {
		return nil, errors.New("at least one spec must be defined")
	}

	config.StartSleep, err = time.ParseDuration(config.StartSleepStr)
	if err != nil {
		return nil, fmt.Errorf("invalid startSleep: %w", err)
	}
	for i, spec := range config.Specs {
		spec.TargetSleep, err = time.ParseDuration(spec.TargetSleepStr)
		if err != nil {
			return nil, fmt.Errorf("invalid spec index %d: %w", i, err)
		}

		spec.Duration, err = time.ParseDuration(spec.DurationStr)
		if err != nil {
			return nil, fmt.Errorf("invalid spec index %d: %w", i, err)
		}

		config.Specs[i] = spec
	}

	return &config, nil
}
