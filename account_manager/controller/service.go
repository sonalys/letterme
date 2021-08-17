package controller

import "github.com/sonalys/letterme/domain/models"

// Service represents the api logic controller,
// It uses of decoupled dependencies to execute business logic for specific cases.
// Examples:
//
// ResetPublicKey: should use persistence layer to fetch account from ownershipID, then update it's publicKey
// to the new one, and delete all pending emails and attachments.
type Service struct {
	*Dependencies
}

// Dependencies are the integrations required to initialize the service.
type Dependencies struct {
}

// NewService initializes the api controller
//
// Here is where you want to start all sub-routines, dependencies validations, etc...
func NewService(d *Dependencies) (*Service, error) {
	return &Service{Dependencies: d}, nil
}

// GetPublicKey gets the publicKey associated with the given address,
// will return error if the address doesn't exist.
func (s *Service) GetPublicKey(address models.Address) (publicKey *models.PublicKey, err error) {
	return nil, nil
}

// ResetPublicKey will re-create the publicKey for the given accountID.
func (s *Service) ResetPublicKey(accountID models.DatabaseID) (account *models.Account, err error) {
	return nil, nil
}

func (s *Service) AddNewDevice(accountID models.DatabaseID) (encryptedPrivateKey *models.EncryptedBuffer, err error) {
	return nil, nil
}

func (s *Service) Authenticate(Address models.Address) (encryptedJWT *models.EncryptedBuffer, err error) {
	return nil, nil
}
