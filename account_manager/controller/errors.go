package controller

import (
	"fmt"

	"github.com/sonalys/letterme/domain/models"
)

const (
	errAddressInUse     = "the address '%s' is already in use"
	errAccountOperation = "failed to %s account: %s"
)

func newAddressInError(address models.Address) error {
	return fmt.Errorf(errAddressInUse, address)
}

func newAccountOperationError(opName string, err error) error {
	return fmt.Errorf(errAccountOperation, opName, err)
}
