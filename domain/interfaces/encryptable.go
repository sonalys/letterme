package interfaces

import "crypto/rsa"

// Encryptable represents objects that can be encrypted.
type Encryptable interface {
	Encrypt(*rsa.PublicKey) error
}
