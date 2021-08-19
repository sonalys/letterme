package controller

import (
	dModels "github.com/sonalys/letterme/domain"
)

// nolint:unused // will be.
func (s *Service) encrypt(k *dModels.PublicKey, src interface{}) (buf *dModels.EncryptedBuffer, err error) {
	return s.CryptographicRouter.Encrypt(k, src)
}

func (s *Service) decrypt(k *dModels.PrivateKey, b *dModels.EncryptedBuffer, dst interface{}) error {
	return s.CryptographicRouter.Decrypt(k, b, dst)
}
