package handler

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sonalys/letterme/domain/messaging/contracts"
	"github.com/sonalys/letterme/domain/models"
	"github.com/sonalys/letterme/domain/persistence/mongo"
)

func (h *Handler) verifyEmailExistence(ctx context.Context, d models.Delivery) (interface{}, error) {
	req := new(contracts.CheckEmailRequest)
	if err := d.GetBody(req); err != nil {
		return nil, errors.Wrap(errDecode, err.Error())
	}

	_, err := h.GetPublicKey(ctx, req.Address)
	switch err {
	case nil:
		return contracts.CheckEmailResponse{Exists: true}, nil
	case mongo.ErrNotFound:
		return contracts.CheckEmailResponse{Exists: false}, nil
	default:
		return nil, errors.Wrap(errInternal, err.Error())
	}
}
