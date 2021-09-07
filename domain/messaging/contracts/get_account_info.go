package contracts

import (
	"github.com/sonalys/letterme/domain/models"
)

type GetAccountInfoRequest struct {
	Address models.Address `json:"address"`
}

type GetAccountInfoResponse struct {
	*models.AccountAddressInfo `json:"account"`
}
