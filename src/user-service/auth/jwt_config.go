package auth

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

type JwtConfig struct {
	PrivateKey string `yaml:"privateKey" env:"PRIVATE_KEY,notEmpty"`
}

func (config JwtConfig) ReadPrivateKey() (any, error) {
	var block *pem.Block

	if _, err := os.Stat(config.PrivateKey); err == nil {
		bytes, err := os.ReadFile(config.PrivateKey)
		if err != nil {
			return nil, err
		}
		block, _ = pem.Decode(bytes)
		if block == nil {
			return nil, errors.New("token or path invalid or empty")
		}
	} else {
		block, _ = pem.Decode([]byte(config.PrivateKey))
		if block == nil {
			return nil, errors.New("token invalid or empty")
		}
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil

}
