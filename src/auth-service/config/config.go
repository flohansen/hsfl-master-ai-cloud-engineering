package config

import (
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/auth-service/auth"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/auth-service/database"
)

type Config struct {
	HttpServerPort string              `env:"HTTP_SERVER_PORT" envDefault:"3000"`
	GrpcServerPort string              `env:"GRPC_SERVER_PORT" envDefault:"50051"`
	Database       database.PsqlConfig `envPrefix:"DB_"`
	JwtConfig      auth.JwtConfig      `envPrefix:"JWT_"`
}
