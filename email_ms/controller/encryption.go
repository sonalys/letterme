package controller

import (
	"github.com/sonalys/letterme/domain/cryptography"
)

func (s *Service) encrypt(k *cryptography.PublicKey, src interface{}) (buf *cryptography.EncryptedBuffer, err error) {
	return s.CryptographicRouter.Encrypt(k, src)
}

func (s *Service) decrypt(k *cryptography.PrivateKey, b *cryptography.EncryptedBuffer, dst interface{}) error {
	return s.CryptographicRouter.Decrypt(k, b, dst)
}
