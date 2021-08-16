package main

import "time"

type DatabaseID string

type Address string

type URL string

// Attachment is a file, image, or any other file type present at the email.
// All external links must be created on our side, one copy for each receiver
// It will be deleted after TTL or readness confirmation.
type Attachment struct {
	// ID is only used by the database system. Do not use it inside code.
	ID        DatabaseID `json:"id"`
	EmailID   string     `json:"email_id"`
	CreatedAt time.Time  `json:"created_at"`
	// ValidUntil is a cached TTL date for the file deletion into the database,
	// different users can have different ttl for their data.
	ValidUntil time.Time `json:"valid_until"`

	// URL is the link for this encrypted attachment on the storage provider.
	URL URL `json:"url"`

	// Size represents the attachment size in bytes.
	Size     uint64 `json:"size"`
	MimeType string `json:"mime_type"`

	// SHA512 is present when the file was sent decrypted, so we can calculate the checksum and verify for vulnerable files.
	SHA512 *string `json:"sha512,omitempty"`
	// Insecure flags attachments that were received without encryption.
	Insecure bool `json:"insecure"`
}

// Email holds information about email and it's metadata.
// It will be deleted after TTL, sent to outside letter.me, or readness confirmation
//
// it will be deleted with all it's attachments.
type Email struct {
	// ID is only used by the database system. Do not use it inside code.
	ID DatabaseID `json:"id"`
	// ValidUntil is a cached TTL date for the file deletion into the database,
	// different users can have different ttl for their data.
	ValidUntil time.Time `json:"valid_until"`

	From        Address      `json:"from"`
	To          []Address    `json:"to"`
	Title       string       `json:"title"`
	Body        string       `json:"body"`
	BodyLength  uint32       `json:"body_length"`
	Attachments []Attachment `json:"attachments,omitempty"`
	// ReceivedAt is only present on incoming emails, outgoing emails to outside letter.me will not have this field.
	ReceivedAt *time.Time `json:"received_at"`
	// SentAt is only present when the email is sent from letter.me.
	SentAt *time.Time `json:"sent_at"`
	// Insecure flags emails that were received without encryption.
	Insecure bool `json:"insecure"`
	// OriginConfirmed is a flag set when the domain from sender confirms the origin of this email.
	OriginConfirmed bool `json:"origin_confirmed"`
}

func (e Email) Validate() error {
	var errMessages []error

	if (len(e.ID) > 0) {

	}

	if len(errMessages) > 0 {
	}
	return nil
}
