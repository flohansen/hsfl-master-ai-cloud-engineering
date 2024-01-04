package auth

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/auth/utils"
	"os"
	"testing"
)

func TestJwtConfig_ReadPrivateKey(t *testing.T) {
	t.Run("Valid private key file", func(t *testing.T) {
		validKeyPath := "./test-key.pem"
		pemString := utils.GenerateRandomECDSAPrivateKeyAsPEM()
		os.WriteFile(validKeyPath, []byte(pemString), 0600)
		defer os.Remove(validKeyPath)

		config := JwtConfig{PrivateKey: validKeyPath}
		_, err := config.ReadPrivateKey()
		if err != nil {
			t.Errorf("Failed to read valid private key file: %v", err)
		}
	})

	t.Run("Invalid private key file", func(t *testing.T) {
		invalidKeyPath := "./test-key.pem"
		os.WriteFile(invalidKeyPath, []byte("invalid key content"), 0600)
		defer os.Remove(invalidKeyPath)

		config := JwtConfig{PrivateKey: invalidKeyPath}
		_, err := config.ReadPrivateKey()
		if err == nil {
			t.Errorf("Expected an error for invalid private key file, but got none")
		}
	})

	t.Run("Valid private key content, no file", func(t *testing.T) {
		config := JwtConfig{PrivateKey: utils.GenerateRandomECDSAPrivateKeyAsPEM()}
		_, err := config.ReadPrivateKey()
		if err != nil {
			t.Errorf("Failed to read valid private key content: %v", err)
		}
	})

	t.Run("Invalid private key content, no file", func(t *testing.T) {
		config := JwtConfig{PrivateKey: "invalid key content"}
		_, err := config.ReadPrivateKey()
		if err == nil {
			t.Errorf("Expected an error for invalid private key content, but got none")
		}
	})
}
