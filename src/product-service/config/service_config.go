package config

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/config/database"
)

type ServiceConfiguration struct {
	Database              database.RQLiteConfig `envPrefix:"RQLITE_"`
	HttpPort              int                   `env:"HTTP_SERVER_PORT" envDefault:"3003"`
	GrpcPort              int                   `env:"GRPC_SERVER_PORT" envDefault:"50053"`
	GrpcUserServiceTarget string                `env:"GRPC_USER_SERVICE_TARGET" envDefault:"localhost:50051"`
}
