package auth

import (
	"crypto/ecdsa"

	"github.com/golang-jwt/jwt"
)

type JwtTokenGenerator struct {
	privateKey *ecdsa.PrivateKey
}

func NewJwtTokenGenerator(config Config) (*JwtTokenGenerator, error) {
	privateKey, err := config.GetPrivateKey()

	if err != nil {
		return nil, err
	}

	return &JwtTokenGenerator{privateKey.(*ecdsa.PrivateKey)}, nil
}

func (g *JwtTokenGenerator) GenerateToken(claims map[string]interface{}) (string, error) {
	jwtClaims := jwt.MapClaims{}

	for key, value := range claims {
		jwtClaims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims)

	return token.SignedString(g.privateKey)
}

func (g *JwtTokenGenerator) ValidateToken(token string) (map[string]interface{}, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return g.privateKey.Public(), nil
	})

	if err != nil {
		return nil, err
	}

	return jwtToken.Claims.(jwt.MapClaims), nil
}
