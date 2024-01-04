package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/auth/utils"
	"strings"
	"testing"
	"time"
)

func TestNewJwtTokenGenerator(t *testing.T) {
	t.Run("Create new JWT Token Generator", func(t *testing.T) {
		t.Run("Invalid Config expect error", func(t *testing.T) {
			var jwtToken, err = NewJwtTokenGenerator(JwtConfig{PrivateKey: "./auth/nonexistent.pem"})
			assert.Nil(t, jwtToken)
			assert.Error(t, err)
		})
		t.Run("Valid Config expect non-nill Token Generator", func(t *testing.T) {
			privateKey := utils.GenerateRandomECDSAPrivateKeyAsPEM()
			tokenGenerator, err := NewJwtTokenGenerator(JwtConfig{PrivateKey: privateKey})
			assert.NoError(t, err)
			assert.NotNil(t, tokenGenerator)
		})
	})
}

func TestJwtTokenGenerator_CreateToken(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	tokenGenerator := JwtTokenGenerator{privateKey}

	token, err := tokenGenerator.CreateToken(map[string]interface{}{
		"exp":  12345,
		"user": "test",
	})

	assert.NoError(t, err)
	tokenParts := strings.Split(token, ".")
	assert.Len(t, tokenParts, 3)

	b, _ := base64.
		StdEncoding.
		WithPadding(base64.NoPadding).
		DecodeString(tokenParts[1])

	var claims map[string]interface{}
	json.Unmarshal(b, &claims)

	assert.Equal(t, float64(12345), claims["exp"])
	assert.Equal(t, "test", claims["user"])
}

func TestJwtTokenGenerator_VerifyToken(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	tokenGenerator := JwtTokenGenerator{privateKey}

	t.Run("should successfully verify JWT token", func(t *testing.T) {
		expiration := time.Now().Add(12345).Unix()
		token, err := tokenGenerator.CreateToken(map[string]interface{}{
			"exp":  float64(expiration),
			"user": "test",
		})

		claims, err := tokenGenerator.VerifyToken(token)

		assert.NoError(t, err)
		assert.Len(t, claims, 2)

		assert.Equal(t, float64(expiration), claims["exp"])
		assert.Equal(t, "test", claims["user"])
	})

	t.Run("should not verify JWT token", func(t *testing.T) {
		token, err := tokenGenerator.CreateToken(map[string]interface{}{
			"exp":  12345,
			"user": "test",
		})

		assert.NoError(t, err)

		_, err = tokenGenerator.VerifyToken(token)
		assert.Error(t, err)
	})
}
