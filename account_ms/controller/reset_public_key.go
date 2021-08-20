package controller

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/sonalys/letterme/account_manager/models"
	dModels "github.com/sonalys/letterme/domain/models"
)

// ResetPublicKey will re-create the publicKey for the given accountID.
func (s *Service) ResetPublicKey(ctx context.Context, req models.ResetPublicKeyRequest) error {
	if err := req.Validate(); err != nil {
		return newInvalidRequestError(err)
	}

	col := s.Persistence.GetCollection(accountCollection)
	queryFilter := dModels.Account{
		OwnershipKey: req.OwnershipKey,
	}
	updateFilter := filter{
		"$set": dModels.Account{
			PublicKey: req.PublicKey,
		},
	}
	if count, err := col.Update(ctx, queryFilter, updateFilter); err != nil {
		return newAccountOperationError("reset", err)
	} else if count != 1 {
		logrus.Errorf("resetPublicKey filtered %d results, should be exactly 1, filter: %#v", count, queryFilter)
	}

	return nil
}
