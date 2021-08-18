package interfaces

import "github.com/sonalys/letterme/domain/cryptography"

type Encryptable interface {
	Encrypt(r cryptography.CryptographicRouter, algorithm cryptography.AlgorithmName, k *cryptography.PublicKey) error
}

type EncryptableValue interface {
	Encrypt(r cryptography.CryptographicRouter, algorithm cryptography.AlgorithmName, k *cryptography.PublicKey) (*cryptography.EncryptedBuffer, error)
}
