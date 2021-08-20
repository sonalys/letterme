package domain

type OwnershipKey string

// EncryptValue implements encryptableValue interface.
func (m OwnershipKey) EncryptValue(r CryptographicRouter, k *PublicKey, algorithm AlgorithmName) (*EncryptedBuffer, error) {
	buf, err := r.EncryptAlgorithm(k, m, algorithm)
	if err == nil {
		return buf, nil
	}
	return nil, newEncryptionError(m, err)
}
