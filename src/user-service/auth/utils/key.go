package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
)

func GenerateRandomECDSAPrivateKeyAsPEM() string {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic("Error generating private key for testing.")
	}

	derFormatKey, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		panic("Error converting ECDSA key to DER format")
	}

	pemKey := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: derFormatKey,
	})

	return string(pemKey)
}
