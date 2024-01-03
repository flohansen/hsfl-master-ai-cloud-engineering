package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJwtAuthorizer(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	tokenGenerator := JwtTokenGenerator{privateKey}

	t.Run("Create new JWT Token Generator", func(t *testing.T) {
		t.Run("Invaid Config expect error", func(t *testing.T) {
			var jwtToken, err = NewJwtTokenGenerator(JwtConfig{PrivateKey: "./auth/test-token-nonexistent"})
			assert.Nil(t, jwtToken)
			assert.Error(t, err)
		})
		t.Run("Valid Config expect non-nill Token Generator", func(t *testing.T) {
			var jwtToken, err = NewJwtTokenGenerator(JwtConfig{PrivateKey: "./test-token"})
			assert.NoError(t, err)
			assert.NotNil(t, jwtToken)
		})
	})

	t.Run("CreateToken", func(t *testing.T) {
		t.Run("should generate valid JWT token", func(t *testing.T) {
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
		})
	})

	t.Run("VerifyToken", func(t *testing.T) {
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
	})

}
