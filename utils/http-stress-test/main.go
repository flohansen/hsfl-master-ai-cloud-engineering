package main

import (
	"http-stress-test/config"
	"http-stress-test/tester"
	"http-stress-test/tester/metrics"
)

func main() {
	cfg := config.GetConfig("config.yaml")

	m := metrics.NewMetrics()
	t := tester.NewTester(cfg, m)

	t.Run()
}
