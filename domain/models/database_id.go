package models

import "errors"

// DatabaseID validates database id's during structural validation.
type DatabaseID string

// Implements interface Validatable.
func (id DatabaseID) Validate() error {
	if len(id) > 0 && len(id) != 40 {
		return newInvalidTypeError("id", string(id), errors.New("length must be 40"))
	}
	return nil
}
