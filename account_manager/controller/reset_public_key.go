package controller

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/sonalys/letterme/domain/cryptography"
	"github.com/sonalys/letterme/domain/models"
)

// ResetPublicKey will re-create the publicKey for the given accountID.
func (s *Service) ResetPublicKey(ctx context.Context, ownershipKey models.OwnershipKey, publicKey cryptography.PublicKey) error {
	col := s.Persistence.GetCollection("account")
	queryFilter := models.Account{
		OwnershipKey: ownershipKey,
	}
	updateFilter := filter{
		"$set": models.Account{
			PublicKey: publicKey,
		},
	}
	if count, err := col.Update(ctx, queryFilter, updateFilter); err != nil {
		return newAccountOperationError("reset", err)
	} else if count != 1 {
		logrus.Errorf("resetPublicKey filtered %d results, should be exactly 1, filter: %#v", count, queryFilter)
	}

	return nil
}
