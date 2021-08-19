package models

import (
	dModels "github.com/sonalys/letterme/domain"
)

// CreateAccountRequest is used to build a request to create a new account.
type CreateAccountRequest struct {
	Address   dModels.Address   `json:"address"`
	PublicKey dModels.PublicKey `json:"public_key"`
}

// Validate implements the validatable interface.
func (r CreateAccountRequest) Validate() error {
	if r.Address == "" {
		return newEmptyFieldError("address")
	}

	if err := r.Address.Validate(); err != nil {
		return newInvalidFieldError("address", err)
	}

	if r.Address.Domain() != "letter.me" {
		return newInvalidFieldError("address", newExternalAddressErr(r.Address))
	}

	if r.PublicKey.IsZero() {
		return newEmptyFieldError("public_key")
	}
	return nil
}
