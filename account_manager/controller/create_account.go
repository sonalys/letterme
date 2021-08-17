package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/sonalys/letterme/domain/models"
)

// CreateAccount receives a new account model, it should be valid, and it's address should not exist already.
func (s *Service) CreateAccount(ctx context.Context, account models.Account) (ownershipToken string, err error) {
	col := s.Persistence.GetCollection("account")
	address := account.Addresses[0]
	var buf models.Account
	if err := col.First(ctx, filter{
		"addresses": filter{"$in": filterList{address}},
	}, &buf); err == nil {
		return "", newAddressInError(address)
	}

	// new accounts should have only 1 address assigned to them.
	account.Addresses = account.Addresses[:1]
	// create new ownership for the account.
	account.OwnershipKey = uuid.NewString()

	if _, err := col.Create(ctx, account); err != nil {
		return "", newAccountOperationError("create", err)
	}

	return account.OwnershipKey, nil
}
