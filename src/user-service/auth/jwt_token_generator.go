package auth

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/golang-jwt/jwt"
)

type JwtTokenGenerator struct {
	privateKey *ecdsa.PrivateKey
}

func NewJwtTokenGenerator(config Config) (*JwtTokenGenerator, error) {
	key, err := config.ReadPrivateKey()
	if err != nil {
		return nil, err
	}

	privateKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, err
	}

	return &JwtTokenGenerator{privateKey}, nil
}

func (gen *JwtTokenGenerator) CreateToken(claims map[string]interface{}) (string, error) {
	jwtClaims := jwt.MapClaims{}
	for k, v := range claims {
		jwtClaims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims)
	return token.SignedString(gen.privateKey)
}

func (gen *JwtTokenGenerator) VerifyToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("signing algorithm unknown: %v", token.Header["alg"])
		}
		return &gen.privateKey.PublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
