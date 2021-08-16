package interfaces

import (
	"github.com/sonalys/letterme/domain/models"
)

// Service creates an interface to mock it outside controller layer.
type Service interface {
	CreateAccount(account models.Account) (ownershipToken string, err error)
	GetPublicKey(address models.Address) (publicKey *models.PublicKey, err error)
	ResetPublicKey(ownershipToken string) (account models.Account, err error)
	AddNewDevice(ownershipToken string) (encryptedPrivateKey models.EncryptedBuffer, err error)
}
