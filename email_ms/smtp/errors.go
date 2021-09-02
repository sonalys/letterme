package smtp

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	errInvalidConfig      = "invalid configuration: %v"
	errInvalidField       = "field '%s' is invalid"
	errInvalidCertificate = "failed to initialize tls certificate"
	errInitializeServer   = "failed to initialize smtp server"
)

func newInvalidConfigErr(errList []error) error {
	return fmt.Errorf(errInvalidConfig, errList)
}

func newInvalidFieldErr(name string) error {
	return fmt.Errorf(errInvalidField, name)
}

func newInvalidCertificateErr(err error) error {
	return errors.Wrap(err, errInvalidCertificate)
}

func newInitializeServerErr(err error) error {
	return errors.Wrap(err, errInitializeServer)
}
