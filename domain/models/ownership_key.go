package models

import "github.com/sonalys/letterme/domain/cryptography"

type OwnershipKey string

// EncryptValue implements encryptableValue interface.
func (m OwnershipKey) EncryptValue(r cryptography.CryptographicRouter, k *cryptography.PublicKey, algorithm cryptography.AlgorithmName) (*cryptography.EncryptedBuffer, error) {
	if buf, err := r.EncryptAlgorithm(k, m, algorithm); err == nil {
		return buf, nil
	} else {
		return nil, newEncryptionError(m, err)
	}
}
