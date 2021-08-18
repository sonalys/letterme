package controller

import (
	"context"

	"github.com/sonalys/letterme/account_manager/interfaces"
	"github.com/sonalys/letterme/domain/cryptography"
	"github.com/sonalys/letterme/domain/models"
)

// Alias for db filters.
type filter map[string]interface{}
type filterList []interface{}

// Service represents the api logic controller,
// It uses of decoupled dependencies to execute business logic for specific cases.
// Examples:
//
// ResetPublicKey: should use persistence layer to fetch account from ownershipID, then update it's publicKey
// to the new one, and delete all pending emails and attachments.
type Service struct {
	Context context.Context
	*Dependencies
}

// Dependencies are the integrations required to initialize the service.
type Dependencies struct {
	interfaces.Persistence
	cryptography.CryptographicRouter
}

// NewService initializes the api controller
//
// Here is where you want to start all sub-routines, dependencies validations, etc...
func NewService(ctx context.Context, d *Dependencies) (*Service, error) {
	return &Service{
		Context:      ctx,
		Dependencies: d,
	}, nil
}

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

// ResetPublicKey will re-create the publicKey for the given accountID.
func (s *Service) ResetPublicKey(ctx context.Context, accountID models.DatabaseID) (account *models.Account, err error) {
	return nil, nil
}

func (s *Service) AddNewDevice(ctx context.Context, accountID models.DatabaseID) (encryptedPrivateKey *cryptography.EncryptedBuffer, err error) {
	return nil, nil
}

func (s *Service) Authenticate(ctx context.Context, Address models.Address) (encryptedJWT *cryptography.EncryptedBuffer, err error) {
	return nil, nil
}
