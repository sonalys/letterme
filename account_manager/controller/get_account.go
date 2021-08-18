package controller

import (
	"context"

	"github.com/sonalys/letterme/domain/models"
)

// GetAccount returns all available information about the owners account.
// returns error if ownership doesn't exist.
func (s *Service) GetAccount(ctx context.Context, ownershipToken models.OwnershipKey) (account *models.Account, err error) {
	if ownershipToken == "" {
		return nil, newInvalidRequestError(newEmptyParamError("ownership_key"))
	}

	col := s.Persistence.GetCollection(accountCollection)
	account = &models.Account{
		OwnershipKey: ownershipToken,
	}

	if err := col.First(ctx, account, account); err != nil {
		return nil, newAccountOperationError("get", err)
	}

	return account, nil
}
