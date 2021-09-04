package controller

import (
	"context"

	"github.com/sonalys/letterme/domain/interfaces"
)

const ServiceConfigurationEnv = "LM_EMAIL_SVC_CONFIG"

// Configuration is used to configure predefined values inside service.
type Configuration struct {
}

// Dependencies are the integrations required to initialize the service.
type Dependencies struct {
	interfaces.Persistence
	interfaces.Messaging
	interfaces.Router
}

// Service represents the api logic controller,
// It uses of decoupled dependencies to execute business logic for specific cases.
type Service struct {
	context context.Context
	config  *Configuration
	*Dependencies
}

// NewService initializes the api controller,
//
// Here is where you want to start all sub-routines, dependencies validations, etc...
func NewService(ctx context.Context, c *Configuration, d *Dependencies) (*Service, error) {
	s := &Service{
		context:      ctx,
		config:       c,
		Dependencies: d,
	}
	return s, nil
}
