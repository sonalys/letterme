package controller

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sonalys/letterme/domain/cryptography"
	"github.com/sonalys/letterme/domain/models"
)

// CreateAccount receives a new account model, it should be valid, and it's address should not exist already.
func (s *Service) CreateAccount(ctx context.Context, account models.Account) (ownershipToken *cryptography.EncryptedBuffer, err error) {
	col := s.Persistence.GetCollection("account")
	address := account.Addresses[0]
	buf := new(models.Account)

	if err := col.First(ctx, filter{
		"addresses": filter{
			"$in": filterList{address},
		},
	}, &buf); err == nil {
		return nil, newAddressInError(address)
	}

	dbAccount := models.Account{
		// new accounts should have only 1 address assigned to them.
		Addresses:   account.Addresses[:1],
		PublicKey:   account.PublicKey,
		DeviceCount: 1,
		// create new ownership for the account.
		OwnershipKey: models.OwnershipKey(uuid.NewString()),
		// 30 days default TTL for media.
		TTL: 30 * 24 * time.Hour,
	}

	if _, err := col.Create(ctx, dbAccount); err != nil {
		return nil, newAccountOperationError("create", err)
	}

	return dbAccount.OwnershipKey.EncryptValue(&account.PublicKey)
}
