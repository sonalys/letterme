package domain

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"hash"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// PublicKey is a type encapsulation of the public key.
type PublicKey rsa.PublicKey

// MarshalJSON implements JSON marshaling interface.
func (k PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(x509.MarshalPKCS1PublicKey(k.Get()))
}

// UnmarshalJSON implements JSON unmarshaling interface.
func (k *PublicKey) UnmarshalJSON(data []byte) error {
	buf := new([]byte)
	if err := json.Unmarshal(data, buf); err != nil {
		return err
	}

	publicKey, err := x509.ParsePKCS1PublicKey(*buf)
	if err != nil {
		return err
	}
	*k = PublicKey(*publicKey)
	return nil
}

// MarshalBSONValue implements BSON marshaling interface.
// PublicKey needs custom encoding because bson doesn't know how to do it.
func (k PublicKey) MarshalBSONValue() (bsontype.Type, []byte, error) {
	bytes := x509.MarshalPKCS1PublicKey(k.Get())
	return bson.MarshalValue(bytes)
}

// UnmarshalBSONValue implements BSON unmarshaling interface.
// PublicKey needs custom encoding because bson doesn't know how to do it.
func (k *PublicKey) UnmarshalBSONValue(dataType bsontype.Type, data []byte) error {
	if dataType == bsontype.Null {
		return nil
	}

	_, buf, ok := (bson.RawValue{Type: dataType, Value: data}).BinaryOK()
	if !ok {
		return fmt.Errorf("failed decoding binary to publicKey from %v", dataType)
	}

	publicKey, err := x509.ParsePKCS1PublicKey(buf)
	if err != nil {
		return err
	}

	*k = PublicKey(*publicKey)
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

// MarshalJSON implements JSON marshaling interface.
func (k PrivateKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(x509.MarshalPKCS1PrivateKey(k.Get()))
}

// UnmarshalJSON implements JSON unmarshaling interface.
func (k *PrivateKey) UnmarshalJSON(data []byte) error {
	buf := new([]byte)
	if err := json.Unmarshal(data, buf); err != nil {
		return err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(*buf)
	if err != nil {
		return err
	}
	*k = PrivateKey(*privateKey)
	return nil
}

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
	pk, err := rsa.GenerateKey(rand.Reader, b)
	if err != nil {
		return nil, err
	}
	privateKey := PrivateKey(*pk)
	return &privateKey, nil
}

type rsa_oaep struct {
	cypher []byte
	hash   hash.Hash
}

// Decrypt uses RSA-OAEP decryption algorithm using sha-512 hash.
// We will use this for authencity checks, we don't decrypt any user content ever.
func (r rsa_oaep) Decrypt(k *PrivateKey, b *EncryptedBuffer, dst interface{}) error {
	msgLen := len(b.Buffer)
	step := k.PublicKey.Size()
	encryptedBytes := b.Buffer
	var decryptedBytes []byte

	for startOffset := 0; startOffset < msgLen; startOffset += step {
		endOffset := startOffset + step
		if endOffset > msgLen {
			endOffset = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(r.hash, rand.Reader, k.Get(), encryptedBytes[startOffset:endOffset], r.cypher)
		if err != nil {
			return err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return json.Unmarshal(decryptedBytes, dst)
}

// Encrypt uses RSA-OAEP encryption algorithm using sha-512 hash.
func (r rsa_oaep) Encrypt(k *PublicKey, src interface{}) (*EncryptedBuffer, error) {
	bytes, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}

	publicKey := k.Get()
	msgLen := len(bytes)
	step := publicKey.Size() - 2*r.hash.Size() - 2
	var encryptedBytes []byte

	for startOffset := 0; startOffset < msgLen; startOffset += step {
		endOffset := startOffset + step
		if endOffset > msgLen {
			endOffset = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(r.hash, rand.Reader, k.Get(), bytes[startOffset:endOffset], r.cypher)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return &EncryptedBuffer{
		Buffer:    encryptedBytes,
		Algorithm: RSA_OAEP,
		Hash:      SHA256,
	}, nil
}
