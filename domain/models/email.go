package models

import (
	"fmt"
	"time"

	"github.com/sonalys/letterme/domain/cryptography"
)

const (
	_         = iota
	KB uint32 = 1 << (10 * iota)
	MB
	GB
)

// UnencryptedEmail is used to receive/send emails from/to outside letter.me,
// It contains decrypted content, and should be encrypted immediately after processing.
//
// It needs to be processed first because it's encrypted using the public key of each recipient.
type UnencryptedEmail struct {
	From        Address             `json:"from"`
	ToList      []Address           `json:"to_list"`
	Attachments []AttachmentRequest `json:"attachments"`
	Inlines     []AttachmentRequest `json:"inlines"`
	Title       []byte              `json:"title"`
	Text        []byte              `json:"text"`
	HTML        []byte              `json:"html"`
}

// NewUnencryptedEmail is used by sync.Pool to pre-allocate envelopes.
func NewUnencryptedEmail() *UnencryptedEmail {
	return &UnencryptedEmail{
		Text:   make([]byte, 0, 1*KB),
		HTML:   make([]byte, 0, 1*MB),
		Title:  make([]byte, 0, 1*KB),
		ToList: make([]Address, 0, 10),
	}
}

// InternalEmailRequest is used to receive a partially encrypted email request from api users,
// It has to have unencrypted sender address, unencrypted toList addresses
//
// It will have to link itself with all it's atachments and transform to Email.
type InternalEmailRequest struct {
	Email                           // heritages all fields from the Encrypted email, expect those redeclared below.
	From        Address             `json:"from"`
	ToList      []Address           `json:"to_list"`
	Attachments []AttachmentRequest `json:"attachments"`
}

// Encrypt implements encryptable interface.
func (m *InternalEmailRequest) Encrypt(r cryptography.CryptographicRouter, k *cryptography.PublicKey, algorithm cryptography.AlgorithmName) error {
	if buf, err := r.EncryptAlgorithm(k, m.From, algorithm); err == nil {
		m.Email.From = buf
	} else {
		return newEncryptionError(m, err)
	}

	if buf, err := r.EncryptAlgorithm(k, m.ToList, algorithm); err == nil {
		m.Email.ToList = buf
	} else {
		return newEncryptionError(m, err)
	}

	for i := range m.Attachments {
		if err := m.Attachments[i].Encrypt(r, algorithm, k); err == nil {
			m.Email.Attachments = append(m.Email.Attachments, m.Attachments[i].Attachment)
		} else {
			return newEncryptionError(m, err)
		}
	}

	return nil
}

// Implements interface Validatable.
func (e InternalEmailRequest) Validate() error {
	var errMessages []error

	if err := e.From.Validate(); err != nil {
		errMessages = append(errMessages, newInvalidTypeError("field", "from", err))
	}

	for i := range e.ToList {
		if err := e.ToList[i].Validate(); err != nil {
			errMessages = append(errMessages, newInvalidTypeError("field", fmt.Sprintf("to[%d]", i), err))
		}
	}

	for i := range e.Attachments {
		if err := e.Attachments[i].Validate(); err != nil {
			errMessages = append(errMessages, newInvalidTypeError("field", fmt.Sprintf("attachments[%d]", i), err))
		}
	}

	if len(errMessages) > 0 {
		return newEmailError(errMessages)
	}
	return nil
}

// Email holds information about email and it's metadata.
// It will be deleted after TTL, sent to outside letter.me, or readness confirmation
//
// This structure is only used for already processed emails, emails from outside are converted to this format after processing.
//
// it will be deleted with all it's attachments.
type Email struct {
	// ID is only used by the database system. Do not use it inside code.
	ID        DatabaseID `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	// ValidUntil is a cached TTL date for the file deletion into the database,
	// different users can have different ttl for their data.
	ValidUntil time.Time `json:"valid_until"`
	// ReadCount is a state used to count how many devices confirmed the receivement of this email,
	// when read_count reaches the user.device_count, it will be deleted will all it's attachments.
	ReadCount uint8 `json:"read_count"`

	// From represents the sender address from this email,
	// It has already been processed and its encrypted.
	From *cryptography.EncryptedBuffer `json:"from"`
	// ToList represents a list of addresses that will receive this email,
	// It has already been processed and its encrypted.
	ToList *cryptography.EncryptedBuffer `json:"to_list"`
	// To represents which person of the list this copy is attributed to.
	// For emails sent to multiple people, each one receive one copy.
	// This field cannot be encrypted because it is a relation.
	//
	// This field is used to fetch pending emails later.
	To Address `json:"to"`

	// Title is the title of the email, cannot be empty,
	// it is encrypted on the device.
	Title *cryptography.EncryptedBuffer `json:"title"`
	// Body represents the body of the email,
	// it is encrypted on the device.
	Body *cryptography.EncryptedBuffer `json:"body"`
	// BodyLength represents the length of the encrypted body chunk.
	BodyLength uint32 `json:"body_length"`
	// Attachments that are already encrypted and hosted inside letter.me.
	Attachments []Attachment `json:"attachments,omitempty"`
	// Insecure flags emails that were received without encryption.
	Insecure bool `json:"insecure"`
	// OriginConfirmed is a flag set when the domain from sender confirms the origin of this email.
	OriginConfirmed bool `json:"origin_confirmed"`
}

// Implements interface Validatable.
func (e Email) Validate() error {
	var errMessages []error

	if err := e.ID.Validate(); err != nil {
		errMessages = append(errMessages, newInvalidTypeError("field", "id", err))
	}

	if len(e.To) == 0 {
		errMessages = append(errMessages, newEmptyFieldError("to"))
	}

	if len(e.Title.Buffer) == 0 {
		errMessages = append(errMessages, newEmptyFieldError("title"))
	}

	if len(e.From.Buffer) == 0 {
		errMessages = append(errMessages, newEmptyFieldError("from"))
	}

	if e.ToList == nil {
		errMessages = append(errMessages, newEmptyFieldError("to_list"))
	}

	if e.BodyLength == 0 && len(e.Attachments) == 0 {
		errMessages = append(errMessages, newEmptyEmailError())
	}

	if len(errMessages) > 0 {
		return newEmailError(errMessages)
	}
	return nil
}
