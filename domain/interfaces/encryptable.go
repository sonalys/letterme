package interfaces

import "github.com/sonalys/letterme/domain/cryptography"

type Encryptable interface {
	Encrypt(r cryptography.CryptographicRouter, k *cryptography.PublicKey, algorithm cryptography.AlgorithmName) error
}

type EncryptableValue interface {
	Encrypt(r cryptography.CryptographicRouter, k *cryptography.PublicKey, algorithm cryptography.AlgorithmName) (*cryptography.EncryptedBuffer, error)
}
