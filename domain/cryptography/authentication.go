package cryptography

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/golang-jwt/jwt"
)

const JWT_AUTH_ENV = "LM_JWT_CONFIG"

// Configuration is the data required to configure a new Auth service.
type AuthConfiguration struct {
	PrivateKey *PrivateKey `json:"private_key"`
}

func (c AuthConfiguration) Validate() error {
	var errList []error
	if c.PrivateKey.E == 0 {
		errList = append(errList, newEmptyFieldError("private_key"))
	}

	if len(errList) > 0 {
		return newInvalidConfigError(c, errList)
	}
	return nil
}

// JWTAuthenticator is a authentication manager, to auth, deauth and authenticate access tokens.
type JWTAuthenticator struct {
	privateKey *PrivateKey
}

func NewJWTAuthenticator(c *AuthConfiguration) *JWTAuthenticator {
	return &JWTAuthenticator{
		privateKey: c.PrivateKey,
	}
}

// Claim is an encapsulation of jwt.Claim.
type Claim jwt.Claims

// CreateToken generates a new token for the given address
func (a *JWTAuthenticator) CreateToken(claims Claim) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims = claims

	buf, err := t.SignedString(a.privateKey.Get())
	if err != nil {
		return "", newOperationJWTError("sign", err)
	}

	return buf, nil
}

// ReadToken parses a token from a buffer and returns it's claims.
func (a *JWTAuthenticator) ReadToken(buf string, dst interface{}) error {
	ref := reflect.TypeOf(dst)
	if ref.Kind() != reflect.Ptr {
		return newInvalidFieldError("dst", fmt.Errorf("%T must be a pointer", dst))
	}

	elem := reflect.ValueOf(dst).Elem()

	if !elem.CanSet() {
		return newInvalidFieldError("dst", fmt.Errorf("cant write to %T", dst))
	}

	addr := reflect.New(elem.Type())
	claim, ok := addr.Interface().(jwt.Claims)
	if !ok {
		return newOperationJWTError("parse", errors.New("not a jwt.Claims interface"))
	}

	token, err := jwt.ParseWithClaims(buf, claim, func(t *jwt.Token) (interface{}, error) {
		return a.privateKey.GetPublicKey().Get(), nil
	})
	if err != nil {
		return newOperationJWTError("parse", err)
	}

	elem.Set(reflect.ValueOf(token.Claims).Elem())
	return nil
}
