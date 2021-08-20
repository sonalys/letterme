package models

import (
	"fmt"

	"github.com/pkg/errors"
	dModels "github.com/sonalys/letterme/domain/models"
)

const (
	errEmptyField      = "field '%s' cannot be empty"
	errInvalidField    = "field '%s' is invalid"
	errExternalAddress = "address '%s' is not from letter.me"
)

func newEmptyFieldError(fieldName string) error {
	return fmt.Errorf(errEmptyField, fieldName)
}

func newInvalidFieldError(fieldName string, err error) error {
	return errors.Wrap(err, fmt.Sprintf(errInvalidField, fieldName))
}

func newExternalAddressErr(address dModels.Address) error {
	return fmt.Errorf(errExternalAddress, address)
}
