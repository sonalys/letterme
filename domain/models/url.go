package models

import "net/url"

type URL string

// Implements interface Validatable.
func (u URL) Validate() error {
	if u == "" {
		return nil
	}

	_, err := url.Parse(string(u))
	if err != nil {
		return newInvalidTypeError("url", string(u), err)
	}

	return nil
}
