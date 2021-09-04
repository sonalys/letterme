package controller

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/sonalys/letterme/domain/interfaces"
	"github.com/sonalys/letterme/domain/persistence/mongo"
	"github.com/sonalys/letterme/domain/utils"
)

const ServiceConfigurationEnv = "LM_EMAIL_SVC_CONFIG"

// Configuration is used to configure predefined values inside service.
type Configuration struct {
}

// Dependencies are the integrations required to initialize the service.
type Dependencies struct {
	interfaces.Persistence
}

// Service represents the api logic controller,
// It uses of decoupled dependencies to execute business logic for specific cases.
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

	svcConfig := new(Configuration)
	if err := utils.LoadFromEnv(ServiceConfigurationEnv, svcConfig); err != nil {
		logrus.Panicf("failed to initialize authConfig from env: %s", err)
	}

	return NewService(ctx, svcConfig, &Dependencies{
		Persistence: mongo,
	})
}
