package contracts

import (
	"github.com/sonalys/letterme/domain/cryptography"
	"github.com/sonalys/letterme/domain/models"
)

type CheckEmailRequest struct {
	Address models.Address `json:"address"`
}

type CheckEmailResponse struct {
	Exists    bool                    `json:"exists"`
	PublicKey *cryptography.PublicKey `json:"public_key"`
}
