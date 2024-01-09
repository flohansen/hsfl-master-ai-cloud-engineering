package config

type Config struct {
	HttpServerPort             string `env:"HTTP_SERVER_PORT" envDefault:"3000"`
	Image                      string `env:"IMAGE"`
	NetworkName                string `env:"NETWORK_NAME"`
	Replicas                   int    `env:"REPLICAS"`
	HealthCheckIntervalSeconds int    `env:"HEALTH_CHECK_INTERVAL_SECONDS"`
}
