package mongo

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("no documents were found")
)

const (
	errInvalidConfiguration = "invalid mongo configuration: %v"
	errInstance             = "failed to create mongo instance"
	errConnect              = "failed to connect to mongo"
	errCustom               = "database error"
	errDecode               = "failed to decode result to dst"
	errCast                 = "failed to cast from %T to %T"
	errOperation            = "error %s documents"
	errCollectionOperation  = "failed to %s collection '%s'"
	errEmptyField           = "field '%s' is empty"
)

func newEmptyFieldError(name string) error {
	return fmt.Errorf(errEmptyField, name)
}

func newInvalidConfigurationError(errList []error) error {
	return fmt.Errorf(errInvalidConfiguration, errList)
}

func newInstanceError(err error) error {
	return errors.Wrap(err, errInstance)
}

func newConnectError(err error) error {
	return errors.Wrap(err, errConnect)
}

func newCustomError(err error) error {
	return errors.Wrap(err, errCustom)
}

func newDecodeError(err error) error {
	return errors.Wrap(err, errDecode)
}

func newCastError(t1, t2 interface{}) error {
	return fmt.Errorf(errCast, t1, t2)
}

func newOperationError(opName string, err error) error {
	return errors.Wrap(err, fmt.Sprintf(errOperation, opName))
}

func newCollectionOperationError(opName, colName string, err error) error {
	return errors.Wrap(err, fmt.Sprintf(errCollectionOperation, opName, colName))
}
