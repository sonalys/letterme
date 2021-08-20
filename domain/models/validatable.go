package models

// Validatable can be used to check if a structure has validation.
type Validatable interface {
	Validate() error
}
