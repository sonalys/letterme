package controller

import (
	"context"

	"github.com/sonalys/letterme/domain"
)

// Authenticate generates a new encrypted jwt token for the given address.
func (s *Service) Authenticate(ctx context.Context, address domain.Address) (encryptedJWT *domain.EncryptedBuffer, err error) {
	if err := address.Validate(); err != nil {
		return nil, newInvalidRequestError(err)
	}

	var publicKey *domain.PublicKey
	publicKey, err = s.GetPublicKey(ctx, address)
	if err != nil {
		return nil, newAccountOperationError("authenticate", err)
	}

	var token string
	token, err = s.Authenticator.CreateToken(address)
	if err != nil {
		return nil, newAccountOperationError("authenticate", err)
	}

	encryptedJWT, err = s.encrypt(publicKey, token)
	if err != nil {
		return nil, newAccountOperationError("authenticate", err)
	}

	return encryptedJWT, nil
}
