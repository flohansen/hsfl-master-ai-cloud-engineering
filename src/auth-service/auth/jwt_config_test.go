package auth

import (
	"crypto/ecdsa"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPrivateKey(t *testing.T) {
	// given
	tmpFile, err := os.CreateTemp("", "privateKey")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(`-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIBi84Hg3uDeepp83gCyw3DQ16bQ4ISl0Af3vvoFbCyhdoAoGCCqGSM49
AwEHoUQDQgAEuUKGJ1ulXgFZ20H1BZDwW65OLXP2WLY6t0v60QDPBtLZ+G1amSx/
Mj7n6KIl35L/+9Z1C2FVG9kzJ+Xh+he1Dw==
-----END EC PRIVATE KEY-----`)

	if err != nil {
		t.Fatal(err)
	}

	err = tmpFile.Close()

	if err != nil {
		t.Fatal(err)
	}

	// when
	config := JwtConfig{PrivateKey: tmpFile.Name()}

	// test
	privateKey, err := config.GetPrivateKey()

	assert.NoError(t, err)
	assert.NotNil(t, privateKey)
	assert.IsType(t, &ecdsa.PrivateKey{}, privateKey)
}
