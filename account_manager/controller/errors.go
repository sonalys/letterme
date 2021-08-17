package controller

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sonalys/letterme/domain/models"
)

const (
	errAddressInUse     = "the address '%s' is already in use"
	errAccountOperation = "failed to %s account"
)

func newAddressInError(address models.Address) error {
	return fmt.Errorf(errAddressInUse, address)
}

func newAccountOperationError(opName string, err error) error {
	return errors.Wrap(err, fmt.Sprintf(errAccountOperation, opName))
}
