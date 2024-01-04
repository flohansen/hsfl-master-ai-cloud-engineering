package main

import (
	"context"
	"http-stress-test/config"
	"http-stress-test/metrics"
	"http-stress-test/tester"
	"log"
)

const configPath = "config.yaml"

func main() {
	// Cancel context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Couldn't load config: %v", err)
	}

	// Metric display
	m := metrics.NewMetrics()
	go m.DisplayMetrics(ctx)

	// Run stress test
	t := tester.NewTester(cfg, m)
	t.Run()

}
