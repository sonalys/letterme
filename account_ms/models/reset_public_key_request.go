package models

import (
	dModels "github.com/sonalys/letterme/domain"
)

// ResetPublicKeyRequest is used to build a request to reset an account.
type ResetPublicKeyRequest struct {
	OwnershipKey dModels.OwnershipKey `json:"ownership_key"`
	PublicKey    dModels.PublicKey    `json:"public_key"`
}

// Validate implements the validatable interface.
func (r ResetPublicKeyRequest) Validate() error {
	if r.OwnershipKey == "" {
		return newEmptyFieldError("ownership_key")
	}

	if r.PublicKey.IsZero() {
		return newEmptyFieldError("public_key")
	}
	return nil
}
