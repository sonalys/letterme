package domain

import (
	"errors"
	"time"
)

// AttachmentRequest is a structure that holds unencrypted information about attachments,
// As soon as the processing phase is done, it will be encrypted into Attachment.
type AttachmentRequest struct {
	Attachment // heritages all fields from the Encrypted email, expect those redeclared below.
	URL        `json:"url"`
}

// Encrypt implements encryptable interface.
func (m AttachmentRequest) Encrypt(r CryptographicRouter, algorithm AlgorithmName, k *PublicKey) error {
	if buf, err := r.EncryptAlgorithm(k, m.URL, algorithm); err == nil {
		m.Attachment.URL = *buf
	} else {
		return newEncryptionError(m, err)
	}
	return nil
}

// Implements interface Validatable.
func (a AttachmentRequest) Validate() error {
	var errMessages []error

	if a.URL == "" {
		errMessages = append(errMessages, newEmptyFieldError("url"))
	} else if err := a.URL.Validate(); err != nil {
		errMessages = append(errMessages, newInvalidTypeError("field", "url", err))
	}

	if len(errMessages) > 0 {
		return newAttachmentError(errMessages)
	}
	return nil
}

// Attachment is a file, image, or any other file type present at the email.
// All external images must be created on our side, one copy for each receiver
// It will be deleted after TTL or readness confirmation.
type Attachment struct {
	// ID is only used by the database system. Do not use it inside code.
	ID        DatabaseID `json:"id"`
	EmailID   DatabaseID `json:"email_id"`
	CreatedAt time.Time  `json:"created_at"`
	// ValidUntil is a cached TTL date for the file deletion into the database,
	// different users can have different ttl for their data.
	ValidUntil time.Time `json:"valid_until"`

	// URL is the link for this encrypted attachment on the storage provider.
	URL EncryptedBuffer `json:"url"`

	// Size represents the attachment size in bytes.
	Size     uint64   `json:"size"`
	MimeType MimeType `json:"mime_type"`

	// SHA512 is present when the file was sent decrypted, so we can calculate the checksum and verify for vulnerable files.
	SHA512 *string `json:"sha512,omitempty"`
	// Insecure flags attachments that were received without encryption.
	Insecure bool `json:"insecure"`
}

// Implements interface Validatable.
func (a Attachment) Validate() error {
	var errMessages []error

	if err := a.ID.Validate(); err != nil {
		errMessages = append(errMessages, newInvalidTypeError("field", "id", err))
	}

	if err := a.EmailID.Validate(); err != nil {
		errMessages = append(errMessages, newInvalidTypeError("field", "email_id", err))
	}

	if a.CreatedAt.IsZero() {
		errMessages = append(errMessages, newEmptyFieldError("created_at"))
	}

	if time.Now().Before(a.ValidUntil) {
		errMessages = append(errMessages, newInvalidTypeError("field", "valid_until", errors.New("date cannot be in the past")))
	}

	if a.Size == 0 {
		errMessages = append(errMessages, newInvalidTypeError("field", "size", errors.New("size cannot be 0")))
	}

	if a.MimeType == "" {
		errMessages = append(errMessages, newEmptyFieldError("mime_type"))
	}

	if err := a.MimeType.Validate(); err != nil {
		errMessages = append(errMessages, newInvalidTypeError("field", "mime_type", err))
	}

	if a.SHA512 != nil && *a.SHA512 == "" {
		errMessages = append(errMessages, newInvalidTypeError("field", "sha512", errors.New("if present, cannot be empty")))
	}

	if len(errMessages) > 0 {
		return newAttachmentError(errMessages)
	}
	return nil
}
