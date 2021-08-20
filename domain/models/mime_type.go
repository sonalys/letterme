package models

import "mime"

type MimeType string

func (m MimeType) Validate() error {
	if m == "" {
		return nil
	}

	_, _, err := mime.ParseMediaType(string(m))
	if err != nil {
		return newInvalidTypeError("mime_type", string(m), err)
	}

	return nil
}
