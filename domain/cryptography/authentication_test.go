package cryptography

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func Test_JWTAuthenticator(t *testing.T) {
	privateKey, err := NewPrivateKey(2048)
	require.NoError(t, err, "should create new private key")

	instance := NewJWTAuthenticator(&AuthConfiguration{
		PrivateKey: privateKey,
	})
	require.NotNil(t, instance, "jwt authenticator instance should be set")

	now := time.Now().Add(time.Hour).Unix()
	expected := jwt.StandardClaims{
		Id:        "123",
		ExpiresAt: now,
	}

	buf, err := instance.CreateToken(expected)
	require.NoError(t, err, "should create instance without errors")

	got := new(jwt.StandardClaims)
	err = instance.ReadToken(buf, got)
	require.NoError(t, err, "should deserialize buffer as decrypted jwt")

	require.EqualValues(t, expected, *got)
}
