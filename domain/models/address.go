package models

import (
	"errors"
	"regexp"
)

type Address string

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Implements interface Validatable.
// does not guarantee that the provided email domain exists.
func (m Address) Validate() error {
	length := len(m)
	if length < 3 || length > 254 {
		return newInvalidTypeError("address", string(m), errors.New("length must be between 4 and 253"))
	}

	if isValid := emailRegex.MatchString(string(m)); !isValid {
		return newInvalidTypeError("address", string(m), errors.New("value is not an email"))
	}
	return nil
}
