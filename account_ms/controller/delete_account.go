package controller

import (
	"context"

	dModels "github.com/sonalys/letterme/domain/models"
)

// DeleteAccount delete the account for the given ownershipToken.
func (s *Service) DeleteAccount(ctx context.Context, ownershipKey dModels.OwnershipKey) (err error) {
	if ownershipKey == "" {
		return newInvalidRequestError(newEmptyParamError("ownership_key"))
	}

	col := s.Persistence.GetCollection(accountCollection)
	if _, err := col.Delete(ctx, dModels.Account{
		OwnershipKey: ownershipKey,
	}); err != nil {
		return newAccountOperationError("delete", err)
	}
	return nil
}
