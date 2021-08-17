package controller

import (
	"context"

	"github.com/sonalys/letterme/domain/models"
)

// GetPublicKey gets the publicKey associated with the given address,
// will return error if the address doesn't exist.
func (s *Service) GetPublicKey(ctx context.Context, address models.Address) (publicKey *models.PublicKey, err error) {
	col := s.Persistence.GetCollection("account")
	account := new(models.Account)
	if err := col.First(ctx, filter{
		"addresses": filter{
			"$in": address,
		},
	}, account); err != nil {
		return nil, newAccountOperationError("get", err)
	}

	return &account.PublicKey, nil
}
