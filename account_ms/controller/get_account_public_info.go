package controller

import (
	"context"

	dModels "github.com/sonalys/letterme/domain/models"
)

// GetAccountPublicInfo gets the publicKey associated with the given address,
// will return error if the address doesn't exist.
func (s *Service) GetAccountPublicInfo(ctx context.Context, address dModels.Address) (info *dModels.AccountAddressInfo, err error) {
	if err := address.Validate(); err != nil {
		return nil, newInvalidRequestError(err)
	}

	col := s.Persistence.GetCollection(accountCollection)

	account := new(dModels.Account)
	if err := col.First(ctx, filter{
		"addresses": address,
	}, account); err != nil {
		return nil, newAccountOperationError("get", err)
	}

	return &dModels.AccountAddressInfo{
		Address:      address,
		PublicKey:    account.PublicKey,
		TTL:          account.TTL,
		MaxEmailSize: account.MaxEmailSize,
		MaxInboxSize: account.MaxInboxSize,
	}, nil
}
