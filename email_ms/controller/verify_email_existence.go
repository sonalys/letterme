package controller

import (
	"context"

	"github.com/sonalys/letterme/domain/messaging"
	"github.com/sonalys/letterme/domain/messaging/contracts"
	"github.com/sonalys/letterme/domain/models"
)

// verifyEmailExistence communicates with account-ms to verify if email exists.
func (s *Service) verifyEmailExistence(ctx context.Context, address models.Address) (bool, error) {
	resp := new(contracts.CheckEmailResponse)
	err := s.Router.Communicate(messaging.QAccountMS, models.Message{
		Type: messaging.ECheckEmail,
		Body: contracts.CheckEmailRequest{Address: address},
	}, resp)
	if err != nil {
		return false, err
	}

	return resp.Exists, nil
}
