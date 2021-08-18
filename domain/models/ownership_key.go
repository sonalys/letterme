package models

import "github.com/sonalys/letterme/domain/cryptography"

type OwnershipKey string

// EncryptValue implements encryptableValue interface.
func (m OwnershipKey) EncryptValue(r cryptography.CryptographicRouter, algorithm cryptography.AlgorithmName, k *cryptography.PublicKey) (*cryptography.EncryptedBuffer, error) {
	if buf, err := r.Encrypt(k, algorithm, m); err == nil {
		return buf, nil
	} else {
		return nil, newEncryptionError(m, err)
	}
}
