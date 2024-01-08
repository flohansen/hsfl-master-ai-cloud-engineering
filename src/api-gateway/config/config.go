package config

type Config struct {
	HttpServerPort          string `env:"HTTP_SERVER_PORT" envDefault:"3000"`
	AuthServiceUrl          string `env:"AUTH_SERVICE_URL"`
	BulletinBoardServiceUrl string `env:"BULLETIN_BOARD_SERVICE_URL"`
	FeedServiceUrl          string `env:"FEED_SERVICE_URL"`
	WebServiceUrl           string `env:"WEB_SERVICE_URL"`
}
