package smtp

import (
	"fmt"
)

var (
	errInvalidConfig = "invalid configuration: %#v"
	errInvalidField  = "field '%s' is invalid"
)

func newInvalidConfigErr(errList []error) error {
	return fmt.Errorf(errInvalidConfig, errList)
}

func newInvalidFieldErr(name string) error {
	return fmt.Errorf(errInvalidField, name)
}
