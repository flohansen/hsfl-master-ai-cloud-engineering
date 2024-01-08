package config

type Config struct {
	HttpServerPort string `env:"HTTP_SERVER_PORT" envDefault:"3000"`
}
