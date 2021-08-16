package models

import "crypto/rsa"

// PublicKey is a type encapsulation of the public key.
type PublicKey rsa.PublicKey

// Get returns a copy of the encapsulated rsa.PublicKey
func (k PublicKey) Get() *rsa.PublicKey {
	key := rsa.PublicKey(k)
	return &key
}
