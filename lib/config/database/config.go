package database

type Config interface {
	GetConnectionString() string
}
