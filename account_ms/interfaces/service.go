package interfaces

import (
	"context"

	"github.com/sonalys/letterme/account_manager/models"

	dModels "github.com/sonalys/letterme/domain"
)

// Service creates an interface to mock it outside controller layer.
type Service interface {
	// User interaction
	Authenticate(ctx context.Context, Address dModels.Address) (encryptedJWT *dModels.EncryptedBuffer, err error)
	AddNewDevice(ctx context.Context, accountID dModels.DatabaseID) (encryptedPrivateKey *dModels.EncryptedBuffer, err error)
	CreateAccount(ctx context.Context, account models.CreateAccountRequest) (encryptedOwnershipToken dModels.EncryptedBuffer, err error)
	// Account Maintenance
	ResetPublicKey(ctx context.Context, req models.ResetPublicKeyRequest) error
	DeleteAccount(ctx context.Context, ownershipToken dModels.OwnershipKey) (err error)
	GetAccount(ctx context.Context, ownershipToken dModels.OwnershipKey) (account dModels.Account, err error)
	// Public Information
	GetPublicKey(ctx context.Context, address dModels.Address) (publicKey *dModels.PublicKey, err error)
}
