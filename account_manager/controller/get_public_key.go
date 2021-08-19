package controller

import (
	"context"

	"github.com/sonalys/letterme/domain"
)

// GetPublicKey gets the publicKey associated with the given address,
// will return error if the address doesn't exist.
func (s *Service) GetPublicKey(ctx context.Context, address domain.Address) (publicKey *domain.PublicKey, err error) {
	if err := address.Validate(); err != nil {
		return nil, newInvalidRequestError(err)
	}

	col := s.Persistence.GetCollection(accountCollection)
	account := new(domain.Account)
	if err := col.First(ctx, filter{
		"addresses": address,
	}, account); err != nil {
		return nil, newAccountOperationError("get", err)
	}

	return &account.PublicKey, nil
}
