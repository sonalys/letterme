package domain

import (
	"fmt"
	"github.com/pkg/errors"
)

// ErrorResponse is the schema of default http errors
type ErrorResponse struct {
	Code   uint   `json:"code"`
	Reason string `json:"reason"`
}

const (
	invalidType = "%s '%s' is not valid"
	emptyField  = "field '%s' cannot be empty"

	emailInvalid = "email is not valid: %#v"
	emailEmpty   = "email cannot be empty"

	attachmentInvalid = "attachment is not valid: %#v"

	encryptionError = "failed to encrypt %T"
	decryptionError = "failed to decrypt %T"

	errOperationJWT = "failed to %s jwt"

	errInvalidConfiguration = "%T configuration is not valid: %#v"
)

func newInvalidConfigError(obj interface{}, errList []error) error {
	return fmt.Errorf(errInvalidConfiguration, obj, errList)
}

func newOperationJWTError(op string, err error) error {
	return errors.Wrap(err, fmt.Sprintf(errOperationJWT, op))
}

func newEncryptionError(obj interface{}, err error) error {
	return errors.Wrap(err, fmt.Sprintf(encryptionError, obj))
}

func newDecryptionError(obj interface{}, err error) error {
	return errors.Wrap(err, fmt.Sprintf(decryptionError, obj))
}

func newInvalidTypeError(typeName, fieldName string, err error) error {
	return errors.Wrap(err, fmt.Sprintf(invalidType, typeName, fieldName))
}

func newEmptyFieldError(fieldName string) error {
	return fmt.Errorf(emptyField, fieldName)
}

func newEmailError(errList []error) error {
	return fmt.Errorf(emailInvalid, errList)
}

func newEmptyEmailError() error {
	return errors.New(emailEmpty)
}

func newAttachmentError(errList []error) error {
	return fmt.Errorf(attachmentInvalid, errList)
}
