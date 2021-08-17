package controller

import (
	"fmt"

	"github.com/sonalys/letterme/domain/models"
)

const (
	errAddressInUse = "the address '%s' is already in use"
)

func newAddressInError(address models.Address) error {
	return fmt.Errorf(errAddressInUse, address)
}
