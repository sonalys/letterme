package interfaces

import (
	"github.com/sonalys/letterme/domain/models"
)

// Service creates an interface to mock it outside controller layer.
type Service interface {
	CreateAccount(account models.Account) (ownershipToken string, err error)
	GetPublicKey(address models.Address) (publicKey *models.PublicKey, err error)
	ResetPublicKey(accountID models.DatabaseID) (account *models.Account, err error)
	AddNewDevice(accountID models.DatabaseID) (encryptedPrivateKey *models.EncryptedBuffer, err error)
	Authenticate(Address models.Address) (encryptedJWT *models.EncryptedBuffer, err error)
}
