package controller

import (
	"context"

	"github.com/sonalys/letterme/account_manager/interfaces"

	"github.com/sonalys/letterme/domain"
)

// Alias for db filters.
type filter map[string]interface{}
type filterList []interface{}

// TODO: maybe put this inside domain.
const accountCollection = "account"

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
	domain.CryptographicRouter
	domain.Authenticator
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

// AddNewDevice todo.
func (s *Service) AddNewDevice(ctx context.Context, accountID domain.DatabaseID) (encryptedPrivateKey *domain.EncryptedBuffer, err error) {
	return nil, nil
}
