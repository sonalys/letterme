package controller

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sonalys/letterme/domain"
)

const (
	errAddressInUse     = "the address '%s' is already in use"
	errAccountOperation = "failed to %s account"
	errInvalidRequest   = "the request is not valid"
	errEmptyParam       = "parameter '%s' cannot be empty"
)

func newAddressInError(address domain.Address) error {
	return fmt.Errorf(errAddressInUse, address)
}

func newEmptyParamError(name string) error {
	return fmt.Errorf(errEmptyParam, name)
}

func newAccountOperationError(opName string, err error) error {
	return errors.Wrap(err, fmt.Sprintf(errAccountOperation, opName))
}

func newInvalidRequestError(err error) error {
	return errors.Wrap(err, errInvalidRequest)
}
