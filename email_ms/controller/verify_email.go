package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sonalys/letterme/domain/messaging"
	"github.com/sonalys/letterme/domain/messaging/contracts"
	"github.com/sonalys/letterme/domain/models"
)

// VerifyEmail communicates with account-ms to verify if email exists.
func (s *Service) VerifyEmail(ctx context.Context, address models.Address) (bool, error) {
	msg := models.Message{
		Type:    messaging.ECheckEmail,
		ReplyTo: uuid.New().String(),
		Body:    contracts.CheckEmailRequest{Address: address},
	}

	if err := s.Messaging.Publish(messaging.QAccountMS, msg); err != nil {
		return false, errors.Wrap(err, "failed to send check-email event")
	}

	m := <-s.Router.WaitResponse(msg.ReplyTo)
	if m.Error != nil {
		return false, m.Error
	}

	resp := new(contracts.CheckEmailResponse)
	if err := m.Message.GetBody(resp); err != nil {
		return false, err
	}

	return resp.Exists, nil
}
