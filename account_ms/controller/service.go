package controller

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sonalys/letterme/domain/interfaces"
	"github.com/sonalys/letterme/domain/persistence/mongo"
	"github.com/sonalys/letterme/domain/utils"

	"github.com/sonalys/letterme/domain/cryptography"
	dModels "github.com/sonalys/letterme/domain/models"
)

// Alias for db filters.
type filter map[string]interface{}
type filterList []interface{}

// TODO: maybe put this inside dModels.
const accountCollection = "account"

const ServiceConfigurationEnv = "LM_ACCOUNT_SVC_CONFIG"

// Configuration is used to configure predefined values inside service.
type Configuration struct {
	AuthTimeout time.Duration `json:"auth_timeout"`
}

// Dependencies are the integrations required to initialize the service.
type Dependencies struct {
	interfaces.Persistence
	cryptography.CryptographicRouter
	cryptography.Authenticator
}

// Service represents the api logic controller,
// It uses of decoupled dependencies to execute business logic for specific cases.
// Examples:
//
// ResetPublicKey: should use persistence layer to fetch account from ownershipID, then update it's publicKey
// to the new one, and delete all pending emails and attachments.
type Service struct {
	context context.Context
	config  *Configuration
	*Dependencies
}

// NewService initializes the api controller
//
// Here is where you want to start all sub-routines, dependencies validations, etc...
func NewService(ctx context.Context, c *Configuration, d *Dependencies) (*Service, error) {
	return &Service{
		context:      ctx,
		config:       c,
		Dependencies: d,
	}, nil
}

// InitializeFromEnv initializes the service from env variables.
func InitializeFromEnv(ctx context.Context) (*Service, error) {
	mongoConfig := new(mongo.Configuration)
	if err := utils.LoadFromEnv(mongo.MongoEnv, mongoConfig); err != nil {
		logrus.Panicf("failed to initialize mongoConfig from env: %s", err)
	}
	mongo, err := mongo.NewClient(ctx, mongoConfig)
	if err != nil {
		panic(err)
	}

	cryptographicConfig := new(cryptography.CryptoConfig)
	if err := utils.LoadFromEnv(cryptography.CRYPTO_CYPHER_ENV, cryptographicConfig); err != nil {
		logrus.Panicf("failed to initialize cryptographicConfig from env: %s", err)
	}

	router, err := cryptography.NewCryptoRouter(cryptographicConfig)
	if err != nil {
		panic(err)
	}

	authConfig := new(cryptography.AuthConfiguration)
	if err := utils.LoadFromEnv(cryptography.JWT_AUTH_ENV, authConfig); err != nil {
		logrus.Panicf("failed to initialize authConfig from env: %s", err)
	}
	auth := cryptography.NewJWTAuthenticator(authConfig)

	svcConfig := new(Configuration)
	if err := utils.LoadFromEnv(ServiceConfigurationEnv, svcConfig); err != nil {
		logrus.Panicf("failed to initialize authConfig from env: %s", err)
	}

	return NewService(ctx, svcConfig, &Dependencies{
		Persistence:         mongo,
		CryptographicRouter: router,
		Authenticator:       auth,
	})
}

// AddNewDevice todo.
func (s *Service) AddNewDevice(ctx context.Context, accountID dModels.DatabaseID) (encryptedPrivateKey *cryptography.EncryptedBuffer, err error) {
	return nil, nil
}
