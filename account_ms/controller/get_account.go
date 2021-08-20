package controller

import (
	"context"

	"github.com/sonalys/letterme/domain"
)

// GetAccount returns all available information about the owners account.
// returns error if ownership doesn't exist.
func (s *Service) GetAccount(ctx context.Context, ownershipToken domain.OwnershipKey) (account *domain.Account, err error) {
	if ownershipToken == "" {
		return nil, newInvalidRequestError(newEmptyParamError("ownership_key"))
	}

	col := s.Persistence.GetCollection(accountCollection)
	account = &domain.Account{
		OwnershipKey: ownershipToken,
	}

	if err := col.First(ctx, account, account); err != nil {
		return nil, newAccountOperationError("get", err)
	}

	return account, nil
}
