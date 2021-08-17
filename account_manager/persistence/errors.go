package persistence

import (
	"errors"
	"fmt"
)

const (
	errInvalidConfiguration = "invalid mongo configuration: %s"
	errInstance             = "failed to create mongo instance: %s"
	errConnect              = "failed to connect to mongo: %s"
	errNotFound             = "no documents were found"
	errCustom               = "database error: %s"
	errDecode               = "failed to decode result to dst: %s"
	errCast                 = "failed to cast from %T to %T"
	errOperation            = "error %s documents: %s"
	errCollectionOperation  = "failed to %s collection '%s': %s"
)

func newInvalidConfigurationError(err error) error {
	return fmt.Errorf(errInvalidConfiguration, err)
}

func newInstanceError(err error) error {
	return fmt.Errorf(errInstance, err)
}

func newConnectError(err error) error {
	return fmt.Errorf(errConnect, err)
}

func newCustomError(err error) error {
	return fmt.Errorf(errCustom, err)
}

func newNotFoundError() error {
	return errors.New(errNotFound)
}

func newDecodeError(err error) error {
	return fmt.Errorf(errDecode, err)
}

func newCastError(t1, t2 interface{}) error {
	return fmt.Errorf(errCast, t1, t2)
}

func newOperationError(opName string, err error) error {
	return fmt.Errorf(errOperation, opName, err)
}

func newCollectionOperationError(opName, colName string, err error) error {
	return fmt.Errorf(errCollectionOperation, opName, colName, err)
}
