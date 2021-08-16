package main

import "fmt"

const (
	baseCode = iota
	invalidEmail
)

const (
	emailEmptyField   = "field '%s' cannot be empty"
	emailInvalidField = "field '%s' is not valid"
)

func newEmptyFieldError(fieldName string) error {
	return fmt.Errorf(emailEmptyField, fieldName)
}
