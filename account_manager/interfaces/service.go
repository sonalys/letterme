package interfaces

import (
	"context"

	"github.com/sonalys/letterme/account_manager/models"
	"github.com/sonalys/letterme/domain/cryptography"
	dModels "github.com/sonalys/letterme/domain/models"
)

// Service creates an interface to mock it outside controller layer.
type Service interface {
	// User interaction
	CreateAccount(ctx context.Context, account models.CreateAccountRequest) (encryptedOwnershipToken cryptography.EncryptedBuffer, err error)
	// AddNewDevice(ctx context.Context, accountID dModels.DatabaseID) (encryptedPrivateKey *cryptography.EncryptedBuffer, err error)
	Authenticate(ctx context.Context, Address dModels.Address) (encryptedJWT *cryptography.EncryptedBuffer, err error)
	// Account Maintenance
	DeleteAccount(ctx context.Context, ownershipToken dModels.OwnershipKey) (err error)
	ResetPublicKey(ctx context.Context, ownershipKey dModels.OwnershipKey, publicKey cryptography.PublicKey) error
	GetAccount(ctx context.Context, ownershipToken dModels.OwnershipKey) (account dModels.Account, err error)
	// Public Information
	GetPublicKey(ctx context.Context, address dModels.Address) (publicKey *cryptography.PublicKey, err error)
}
