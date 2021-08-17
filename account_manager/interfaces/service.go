package interfaces

import (
	"context"

	"github.com/sonalys/letterme/domain/models"
)

// Service creates an interface to mock it outside controller layer.
type Service interface {
	CreateAccount(ctx context.Context, account models.Account) (encryptedOwnershipToken models.EncryptedBuffer, err error)
	DeleteAccount(ctx context.Context, ownershipToken models.OwnershipKey) (err error)
	GetAccount(ctx context.Context, ownershipToken models.OwnershipKey) (account models.Account, err error)
	GetPublicKey(ctx context.Context, address models.Address) (publicKey *models.PublicKey, err error)
	ResetPublicKey(ctx context.Context, accountID models.DatabaseID) (account *models.Account, err error)
	AddNewDevice(ctx context.Context, accountID models.DatabaseID) (encryptedPrivateKey *models.EncryptedBuffer, err error)
	Authenticate(ctx context.Context, Address models.Address) (encryptedJWT *models.EncryptedBuffer, err error)
}
