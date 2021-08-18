package cryptography

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"hash"
)

const CRYPTO_CYPHER_ENV = "LM_CRYPTO_CYPHER"

// PublicKey is a type encapsulation of the public key.
type PublicKey rsa.PublicKey

// IsZero implements the empty struct interface.
func (k PublicKey) IsZero() bool {
	return k.Get().N == nil
}

// Get returns a copy of the encapsulated rsa.PublicKey
func (k PublicKey) Get() *rsa.PublicKey {
	key := rsa.PublicKey(k)
	return &key
}

// PrivateKey is a type encapsulation of the private key.
type PrivateKey rsa.PrivateKey

// Get returns a copy of the encapsulated rsa.PublicKey
func (k PrivateKey) Get() *rsa.PrivateKey {
	key := rsa.PrivateKey(k)
	return &key
}

// NewPrivateKey generates a new RSA privateKey.
func NewPrivateKey(b int) (*PrivateKey, error) {
	if pk, err := rsa.GenerateKey(rand.Reader, b); err != nil {
		return nil, err
	} else {
		pk := PrivateKey(*pk)
		return &pk, nil
	}
}

type rsa_oaep struct {
	cypher []byte
	hash   hash.Hash
}

// Decrypt uses RSA-OAEP decryption algorithm using sha-512 hash.
// We will use this for authencity checks, we don't decrypt any user content ever.
func (r rsa_oaep) Decrypt(k *PrivateKey, b *EncryptedBuffer, dst interface{}) error {
	buf, err := rsa.DecryptOAEP(r.hash, rand.Reader, k.Get(), b.Buffer, r.cypher)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, dst)
}

// Encrypt uses RSA-OAEP encryption algorithm using sha-512 hash.
func (r rsa_oaep) Encrypt(k *PublicKey, src interface{}) (*EncryptedBuffer, error) {
	bytes, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}

	encryptedBuf, err := rsa.EncryptOAEP(r.hash, rand.Reader, k.Get(), bytes, r.cypher)
	if err != nil {
		return nil, err
	}

	return &EncryptedBuffer{
		Buffer:    encryptedBuf,
		Algorithm: RSA_OAEP,
		Hash:      sha256,
	}, nil
}
