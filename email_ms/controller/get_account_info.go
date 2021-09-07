package controller

import (
	"context"

	"github.com/sonalys/letterme/domain/messaging"
	"github.com/sonalys/letterme/domain/messaging/contracts"
	"github.com/sonalys/letterme/domain/models"
)

// getAccountInfo communicates with account-ms to fetch email info.
func (s *Service) getAccountInfo(ctx context.Context, address models.Address) (resp *contracts.GetAccountInfoResponse, err error) {
	resp = new(contracts.GetAccountInfoResponse)
	err = s.EventRouter.Communicate(messaging.AccountMS, messaging.Message{
		Type: messaging.FetchEmailPublicInfo,
		Body: contracts.GetAccountInfoRequest{Address: address},
	}, resp)
	if err != nil {
		return nil, err
	}

	return
}
