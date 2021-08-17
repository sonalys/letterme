package controller

import (
	"context"

	"github.com/sonalys/letterme/domain/models"
)

// DeleteAccount delete the account for the given ownershipToken.
func (s *Service) DeleteAccount(ctx context.Context, ownershipToken string) (err error) {
	col := s.Persistence.GetCollection("account")
	if err := col.Delete(ctx, models.Account{
		OwnershipKey: ownershipToken,
	}); err != nil {
		return newAccountOperationError("delete", err)
	}
	return nil
}
