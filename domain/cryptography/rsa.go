package cryptography

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"hash"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// PublicKey is a type encapsulation of the public key.
type PublicKey rsa.PublicKey

// MarshalBSONValue implements BSON marshaling interface.
// PublicKey needs custom encoding because bson doesn't know how to do it.
func (s PublicKey) MarshalBSONValue() (bsontype.Type, []byte, error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return bson.TypeBinary, nil, err
	}
	return bson.MarshalValue(bytes)
}

// UnmarshalBSONValue implements BSON unmarshaling interface.
// PublicKey needs custom encoding because bson doesn't know how to do it.
func (s *PublicKey) UnmarshalBSONValue(dataType bsontype.Type, data []byte) error {
	if dataType == bsontype.Null {
		return nil
	}

	_, buf, ok := (bson.RawValue{Type: dataType, Value: data}).BinaryOK()
	if !ok {
		return fmt.Errorf("failed decode binary to publicKey from %v", dataType)
	}

	var pk PublicKey
	if err := json.Unmarshal([]byte(buf), &pk); err != nil {
		return fmt.Errorf("failed decode object publicKey from %v", dataType)
	}

	*s = pk
	return nil
}

// IsZero implements the empty struct interface.
func (k PublicKey) IsZero() bool {
	return k.Get().E == 0
}

// Get returns a copy of the encapsulated rsa.PublicKey
func (k PublicKey) Get() *rsa.PublicKey {
	key := rsa.PublicKey(k)
	return &key
}

// PrivateKey is a type encapsulation of the private key.
type PrivateKey rsa.PrivateKey

// GetPublicKey returns a copy of the encapsulated rsa.PublicKey
func (k PrivateKey) GetPublicKey() *PublicKey {
	key := PublicKey(k.PublicKey)
	return &key
}

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
