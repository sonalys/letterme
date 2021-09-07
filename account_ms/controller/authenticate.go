package controller

import (
	"context"
	"time"

	"github.com/sonalys/letterme/domain/cryptography"
	dModels "github.com/sonalys/letterme/domain/models"
)

// Authenticate generates a new encrypted jwt token for the given address.
func (s *Service) Authenticate(ctx context.Context, address dModels.Address) (encryptedJWT *cryptography.EncryptedBuffer, err error) {
	if err := address.Validate(); err != nil {
		return nil, newInvalidRequestError(err)
	}

	account, err := s.GetAccountPublicInfo(ctx, address)
	if err != nil {
		return nil, newAccountOperationError("authenticate", err)
	}

	var token string
	token, err = s.Authenticator.CreateToken(&dModels.TokenClaims{
		Address:   address,
		ExpiresAt: time.Now().Add(s.config.AuthTimeout).Unix(),
	})
	if err != nil {
		return nil, newAccountOperationError("authenticate", err)
	}

	encryptedJWT, err = s.encrypt(account.PublicKey, token)
	if err != nil {
		return nil, newAccountOperationError("authenticate", err)
	}

	return encryptedJWT, nil
}
