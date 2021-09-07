package handler

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sonalys/letterme/domain/messaging"
	"github.com/sonalys/letterme/domain/messaging/contracts"
	"github.com/sonalys/letterme/domain/persistence/mongo"
)

func (h *Handler) getAccountPublicInfo(ctx context.Context, d messaging.Delivery) (interface{}, error) {
	req := new(contracts.GetAccountInfoRequest)
	if err := d.GetBody(req); err != nil {
		return nil, errors.Wrap(errDecode, err.Error())
	}

	info, err := h.GetAccountPublicInfo(ctx, req.Address)
	switch err {
	case nil:
		return contracts.GetAccountInfoResponse{AccountAddressInfo: info}, nil
	case mongo.ErrNotFound:
		return contracts.GetAccountInfoResponse{}, nil
	default:
		return nil, errors.Wrap(errInternal, err.Error())
	}
}
