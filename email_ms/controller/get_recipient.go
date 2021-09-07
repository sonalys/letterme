package controller

import (
	"context"

	"github.com/sonalys/letterme/domain/cryptography"
	"github.com/sonalys/letterme/domain/messaging"
	"github.com/sonalys/letterme/domain/messaging/contracts"
	"github.com/sonalys/letterme/domain/models"
)

// getRecipient communicates with account-ms to fetch email info.
func (s *Service) getRecipient(ctx context.Context, address models.Address) (exists bool, pk *cryptography.PublicKey, err error) {
	resp := new(contracts.CheckEmailResponse)
	err = s.EventRouter.Communicate(messaging.AccountMS, messaging.Message{
		Type: messaging.FetchEmailPublicInfo,
		Body: contracts.CheckEmailRequest{Address: address},
	}, resp)
	if err != nil {
		return false, nil, err
	}

	return resp.Exists, resp.PublicKey, nil
}
