package database

import (
	"fmt"
)

type RQLiteConfig struct {
	Host     string `env:"HOST,notEmpty"`
	Port     int    `env:"PORT,notEmpty"`
	Username string `env:"USER,notEmpty"`
	Password string `env:"PASSWORD,notEmpty"`
}

func (r *RQLiteConfig) GetConnectionString() string {
	return fmt.Sprintf("http://%s:%s@%s:%d/",
		r.Host, r.Username, r.Password, r.Port)
}
