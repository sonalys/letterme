package controller

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sonalys/letterme/account_ms/models"
	"github.com/sonalys/letterme/domain/cryptography"
	dModels "github.com/sonalys/letterme/domain/models"
)

// CreateAccount receives a new account model, it should be valid, and it's address should not exist already.
func (s *Service) CreateAccount(ctx context.Context, account models.CreateAccountRequest) (ownershipToken *cryptography.EncryptedBuffer, err error) {
	if err := account.Validate(); err != nil {
		return nil, newInvalidRequestError(err)
	}

	col := s.Persistence.GetCollection(accountCollection)

	if err := col.First(ctx, filter{
		"addresses": filter{
			"$in": filterList{account.Address},
		},
		// err == nil because it will return errNotFound if email is available
	}, &dModels.Account{}); err == nil {
		return nil, newAddressInError(account.Address)
	}

	dbAccount := dModels.Account{
		// new accounts should have only 1 address assigned to them.
		Addresses:   []dModels.Address{account.Address},
		PublicKey:   account.PublicKey,
		DeviceCount: 1,
		// create new ownership for the account.
		OwnershipKey: dModels.OwnershipKey(uuid.NewString()),
		// 30 days default TTL for media.
		TTL: 30 * 24 * time.Hour,
	}

	if _, err := col.Create(ctx, dbAccount); err != nil {
		return nil, newAccountOperationError("create", err)
	}

	return dbAccount.OwnershipKey.EncryptValue(s.CryptographicRouter, &account.PublicKey, cryptography.RSA_OAEP)
}
