package auth

import (
	"github.com/golang-jwt/jwt"
	"os"
)

type JwtConfig struct {
	PrivateKey string `yaml:"privateKey"`
	PublicKey  string `yaml:"publicKey"`
}

func (config JwtConfig) ReadPrivateKey() (any, error) {
	bytes, err := os.ReadFile(config.PrivateKey)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (config JwtConfig) ReadPublicKey() (any, error) {

	bytes, err := os.ReadFile(config.PublicKey)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
