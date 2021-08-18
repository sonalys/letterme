package controller

import "github.com/sonalys/letterme/domain/cryptography"

func (s *Service) encrypt(k *cryptography.PublicKey, src interface{}) (*cryptography.EncryptedBuffer, error) {
	if buf, err := s.CryptographicRouter.Encrypt(k, src); err != nil {
		return nil, err
	} else {
		return buf, nil
	}
}

func (s *Service) decrypt(k *cryptography.PrivateKey, b *cryptography.EncryptedBuffer, dst interface{}) error {
	if err := s.CryptographicRouter.Decrypt(k, b, dst); err != nil {
		return err
	}
	return nil
}
