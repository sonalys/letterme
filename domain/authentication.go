package domain

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const JWT_AUTH_ENV = "LM_JWT_CONFIG"

// Configuration is the data required to configure a new Auth service.
type AuthConfiguration struct {
	PrivateKey     PrivateKey    `json:"private_key"`
	ExpiryDuration time.Duration `json:"expiry_duration"`
}

// TokenClaims are the letter.me customized jwt token claims.
// Need to have address and expiry date.
type TokenClaims struct {
	*jwt.StandardClaims
	Address Address
}

// JWTAuthenticator is a authentication manager, to auth, deauth and authenticate access tokens.
type JWTAuthenticator struct {
	privateKey     PrivateKey
	expiryDuration time.Duration
}

func NewJWTAuthenticator(c *AuthConfiguration) *JWTAuthenticator {
	return &JWTAuthenticator{
		privateKey:     c.PrivateKey,
		expiryDuration: c.ExpiryDuration,
	}
}

// CreateToken generates a new token for the given address
func (a *JWTAuthenticator) CreateToken(address Address) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims = TokenClaims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.expiryDuration).Unix(),
		},
		Address: address,
	}

	buf, err := t.SignedString(a.privateKey.Get())
	if err != nil {
		return "", newOperationJWTError("sign", err)
	}

	return buf, nil
}

// ReadToken parses a token from a buffer and returns it's claims.
func (a *JWTAuthenticator) ReadToken(buf string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(buf, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return a.privateKey.GetPublicKey().Get(), nil
	})
	if err != nil {
		return nil, newOperationJWTError("parse", err)
	}

	if claims, ok := token.Claims.(*TokenClaims); ok {
		return claims, nil
	} else {
		return nil, newOperationJWTError("read", errors.New("failed to parse decoded jwt token to claims"))
	}
}
