package auth

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt"
)

type JwtTokenGenerator struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJwtTokenGenerator(config Config) (*JwtTokenGenerator, error) {
	uncheckedPrivateKey, err := config.ReadPrivateKey()
	if err != nil {
		return nil, err
	}
	uncheckedPublicKey, err := config.ReadPublicKey()
	if err != nil {
		return nil, err
	}

	privateKey, ok := uncheckedPrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, err
	}
	publicKey, ok := uncheckedPublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, err
	}

	return &JwtTokenGenerator{privateKey, publicKey}, nil
}

func (gen *JwtTokenGenerator) CreateToken(claims map[string]interface{}) (string, error) {
	jwtClaims := jwt.MapClaims{}
	for k, v := range claims {
		jwtClaims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtClaims)
	return token.SignedString(gen.privateKey)
}
