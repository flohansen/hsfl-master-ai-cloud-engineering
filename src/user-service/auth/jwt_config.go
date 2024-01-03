package auth

import (
	"crypto/x509"
	"encoding/pem"
	"os"
)

type JwtConfig struct {
	PrivateKey string `yaml:"privateKey"`
}

func (config JwtConfig) ReadPrivateKey() (any, error) {
	var block *pem.Block

	if _, err := os.Stat(config.PrivateKey); err == nil {
		bytes, err := os.ReadFile(config.PrivateKey)
		if err != nil {
			return nil, err
		}
		block, _ = pem.Decode(bytes)
	} else {
		block, _ = pem.Decode([]byte(config.PrivateKey))
		if block == nil {
			return nil, err
		}
	}

	return x509.ParseECPrivateKey(block.Bytes)

}
