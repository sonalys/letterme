package models

import "crypto/rsa"

// PublicKey is a type encapsulation of the public key.
type PublicKey rsa.PublicKey

// IsZero implements the empty struct interface.
func (k PublicKey) IsZero() bool {
	return k.Get().N == nil
}

// Get returns a copy of the encapsulated rsa.PublicKey
func (k PublicKey) Get() *rsa.PublicKey {
	key := rsa.PublicKey(k)
	return &key
}
