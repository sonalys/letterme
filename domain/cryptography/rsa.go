package cryptography

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

var cypher = []byte("6 _$a2&B¨2e3Rk(*789")

// Decrypt uses RSA-OAEP decryption algorithm using sha-512 hash.
// We will use this for authencity checks, we don't decrypt any user content ever.
func Decrypt(k *rsa.PrivateKey, b []byte) ([]byte, error) {
	return rsa.DecryptOAEP(crypto.SHA256.New(), rand.Reader, k, b, cypher)
}

// Encrypt uses RSA-OAEP encryption algorithm using sha-512 hash.
func Encrypt(k *rsa.PublicKey, b []byte) ([]byte, error) {
	return rsa.EncryptOAEP(crypto.SHA256.New(), rand.Reader, k, b, cypher)
}
