package models

import (
	"errors"
	"regexp"
	"strings"
)

// Address is used to reference all email addresses in letter.me.
type Address string

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Implements interface Validatable.
// does not guarantee that the provided email domain exists.
func (m Address) Validate() error {
	length := len(m)
	if length < 4 || length > 253 {
		return newInvalidTypeError("address", string(m), errors.New("length must be between 4 and 253"))
	}

	if isValid := emailRegex.MatchString(string(m)); !isValid {
		return newInvalidTypeError("address", string(m), errors.New("value is not valid"))
	}
	return nil
}

// Domain is used to get the email dModels.
//
// Example:
//	"alysson@letter.me" => "letter.me"
func (m Address) Domain() string {
	if buf := strings.Split(string(m), "@"); len(buf) == 2 {
		return buf[1]
	}
	return ""
}
