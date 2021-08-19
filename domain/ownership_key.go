package domain

type OwnershipKey string

// EncryptValue implements encryptableValue interface.
func (m OwnershipKey) EncryptValue(r CryptographicRouter, k *PublicKey, algorithm AlgorithmName) (*EncryptedBuffer, error) {
	if buf, err := r.EncryptAlgorithm(k, m, algorithm); err == nil {
		return buf, nil
	} else {
		return nil, newEncryptionError(m, err)
	}
}
