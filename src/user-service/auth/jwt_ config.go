package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtConfig struct {
	PrivateKey string `env:"PRIVATE_KEY,notEmpty"`
	PublicKey  string `env:"PUBLIC_KEY,notEmpty"`
}

func (config JwtConfig) ReadPrivateKey() (any, error) {
	bytes := []byte(config.PrivateKey)
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (config JwtConfig) ReadPublicKey() (any, error) {
	bytes := []byte(config.PublicKey)
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
