package models

import (
	"errors"
	"fmt"
)

// ErrorResponse is the schema of default http errors
type ErrorResponse struct {
	Code   uint   `json:"code"`
	Reason string `json:"reason"`
}

const (
	invalidType = "%s '%s' is not valid: %s"
	emptyField  = "field '%s' cannot be empty"

	emailInvalid = "email is not valid: %#v"
	emailEmpty   = "email cannot be empty"

	attachmentInvalid = "attachment is not valid: %#v"

	encryptionError = "failed to encrypt %T: %s"
	decryptionError = "failed to decrypt %T: %s"
)

func newEncryptionError(obj interface{}, err error) error {
	return fmt.Errorf(encryptionError, obj, err)
}

func newDecryptionError(obj interface{}, err error) error {
	return fmt.Errorf(decryptionError, obj, err)
}

func newInvalidTypeError(typeName, fieldName string, err error) error {
	return fmt.Errorf(invalidType, typeName, fieldName, err)
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
