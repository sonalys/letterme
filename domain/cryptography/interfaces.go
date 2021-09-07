package cryptography

// EncryptionAlgorithm will be used to allow our system to use multiple
// encryption algorithms, that can be changed at any time.
type EncryptionAlgorithm interface {
	Decrypt(k *PrivateKey, b *EncryptedBuffer, dst interface{}, cypher []byte) error
	Encrypt(k *PublicKey, src interface{}, cypher []byte) (*EncryptedBuffer, error)
}

// CryptographicRouter is used to encrypt and decrypt EncryptedBuffer struct,
// it is capable of routing through multiple cryptographic algorithms.
type CryptographicRouter interface {
	Decrypt(k *PrivateKey, b *EncryptedBuffer, dst interface{}) error
	Encrypt(k *PublicKey, src interface{}) (*EncryptedBuffer, error)
	EncryptWithCypher(k *PublicKey, src interface{}, cypher []byte) (*EncryptedBuffer, error)
	DecryptWithCypher(k *PublicKey, src interface{}, cypher []byte) (*EncryptedBuffer, error)
	EncryptAlgorithm(k *PublicKey, src interface{}, algorithm AlgorithmName) (*EncryptedBuffer, error)
}

// Authenticator is an entity capable of transforming claims to an encoded string,
// and later decrypt this encrypted buffer to claims again.
type Authenticator interface {
	CreateToken(claim Claim) (string, error)
	ReadToken(buf string, dst interface{}) error
}
