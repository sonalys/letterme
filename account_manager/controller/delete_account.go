package controller

import (
	"context"

	"github.com/sonalys/letterme/domain"
)

// DeleteAccount delete the account for the given ownershipToken.
func (s *Service) DeleteAccount(ctx context.Context, ownershipKey domain.OwnershipKey) (err error) {
	if ownershipKey == "" {
		return newInvalidRequestError(newEmptyParamError("ownership_key"))
	}

	col := s.Persistence.GetCollection(accountCollection)
	if _, err := col.Delete(ctx, domain.Account{
		OwnershipKey: ownershipKey,
	}); err != nil {
		return newAccountOperationError("delete", err)
	}
	return nil
}
