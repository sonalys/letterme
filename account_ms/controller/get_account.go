package controller

import (
	"context"

	dModels "github.com/sonalys/letterme/domain/models"
)

// GetAccount returns all available information about the owners account.
// returns error if ownership doesn't exist.
func (s *Service) GetAccount(ctx context.Context, ownershipToken dModels.OwnershipKey) (account *dModels.Account, err error) {
	if ownershipToken == "" {
		return nil, newInvalidRequestError(newEmptyParamError("ownership_key"))
	}

	col := s.Persistence.GetCollection(accountCollection)
	account = &dModels.Account{
		OwnershipKey: ownershipToken,
	}

	if err := col.First(ctx, account, account); err != nil {
		return nil, newAccountOperationError("get", err)
	}

	return account, nil
}
