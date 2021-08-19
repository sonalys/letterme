package domain

type Encryptable interface {
	Encrypt(r CryptographicRouter, k *PublicKey, algorithm AlgorithmName) error
}

type EncryptableValue interface {
	Encrypt(r CryptographicRouter, k *PublicKey, algorithm AlgorithmName) (*EncryptedBuffer, error)
}
