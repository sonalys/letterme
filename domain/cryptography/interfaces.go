package cryptography

import "hash"

// EncryptionAlgorithm will be used to allow our system to use multiple
// encryption algorithms, that can be changed at any time.
type EncryptionAlgorithm interface {
	Decrypt(k *PrivateKey, b *EncryptedBuffer, dst interface{}) error
	Encrypt(k *PublicKey, src interface{}) (*EncryptedBuffer, error)
}

// CryptographicRouter is used to encrypt and decrypt EncryptedBuffer struct,
// it is capable of routing through multiple cryptographic algorithms.
type CryptographicRouter interface {
	AddRSA_OAEP(cypher []byte, hash hash.Hash)
	Decrypt(k *PrivateKey, b *EncryptedBuffer, dst interface{}) error
	Encrypt(k *PublicKey, algorithm AlgorithmName, src interface{}) (*EncryptedBuffer, error)
}
