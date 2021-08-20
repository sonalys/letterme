package models

import "github.com/sonalys/letterme/domain/cryptography"

type OwnershipKey string

// EncryptValue implements encryptableValue interface.
func (m OwnershipKey) EncryptValue(r cryptography.CryptographicRouter, k *cryptography.PublicKey, algorithm cryptography.AlgorithmName) (*cryptography.EncryptedBuffer, error) {
	buf, err := r.EncryptAlgorithm(k, m, algorithm)
	if err == nil {
		return buf, nil
	}
	return nil, newEncryptionError(m, err)
}
