package models

import "fmt"

const (
	errEmptyField = "field '%s' cannot be empty"
)

func newEmptyFieldError(fieldName string) error {
	return fmt.Errorf(errEmptyField, fieldName)
}
