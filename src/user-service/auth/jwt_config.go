package auth

import (
	"crypto/x509"
	"encoding/pem"
	"os"
)

type JwtConfig struct {
	PrivateKey string `yaml:"signKey"`
	PublicKey  string `yaml:"publicKey"`
}

func (config JwtConfig) ReadPrivateKey() (any, error) {
	bytes, err := os.ReadFile(config.PrivateKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)
	return x509.ParseECPrivateKey(block.Bytes)
}
