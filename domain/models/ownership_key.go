package models

import "github.com/sonalys/letterme/domain/cryptography"

type OwnershipKey string

// Encrypt implements encryptableValue interface.
func (m OwnershipKey) EncryptValue(k *PublicKey) (*EncryptedBuffer, error) {
	if buf, err := cryptography.Encrypt(k.Get(), []byte(m)); err == nil {
		return &EncryptedBuffer{
			Buffer:    buf,
			Algorithm: RSA_OAEP,
		}, nil
	} else {
		return nil, newEncryptionError(m, err)
	}
}
