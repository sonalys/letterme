package cryptography

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

// Decrypt uses RSA-OAEP decryption algorithm using sha-512 hash.
// We will use this for authencity checks, we don't decrypt any user content ever.
func Decrypt(k *rsa.PrivateKey, b []byte) ([]byte, error) {
	return k.Decrypt(nil, b, rsa.OAEPOptions{Hash: crypto.SHA512})
}

// Encrypt uses RSA-OAEP encryption algorithm using sha-512 hash.
func Encrypt(k *rsa.PublicKey, b []byte) ([]byte, error) {
	return rsa.EncryptOAEP(crypto.SHA512.New(), rand.Reader, k, b, nil)
}
