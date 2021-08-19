package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/sonalys/letterme/account_manager/controller"
	"github.com/sonalys/letterme/account_manager/persistence"
	"github.com/sonalys/letterme/account_manager/utils"
	"github.com/sonalys/letterme/domain"
)

func main() {
	// Context with cancel so we can stop all children from their inner loops after os.Interrupt.
	ctx, cancel := context.WithCancel(context.Background())

	// INFO: stop channel is needed for graceful shutdown implementation.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	dep := initializeDependencies(ctx)
	if _, err := controller.NewService(ctx, dep); err != nil {
		panic(err)
	}

	<-stop
	cancel()
}

func initializeDependencies(ctx context.Context) *controller.Dependencies {
	mongoConfig := new(persistence.Configuration)
	if err := utils.LoadFromEnv(persistence.MongoEnv, mongoConfig); err != nil {
		panic("failed to initialize mongoConfig from env")
	}
	mongo, err := persistence.NewMongo(ctx, mongoConfig)
	if err != nil {
		panic(err)
	}

	cryptographicConfig := new(domain.CryptoConfig)
	if err := utils.LoadFromEnv(domain.CRYPTO_CYPHER_ENV, cryptographicConfig); err != nil {
		panic("failed to initialize cryptographicConfig from env")
	}

	router, err := domain.NewCryptoRouter(cryptographicConfig)
	if err != nil {
		panic(err)
	}

	authConfig := new(domain.AuthConfiguration)
	if err := utils.LoadFromEnv(domain.JWT_AUTH_ENV, authConfig); err != nil {
		panic("failed to initialize authConfig from env")
	}

	auth := domain.NewJWTAuthenticator(authConfig)

	return &controller.Dependencies{
		Persistence:         mongo,
		CryptographicRouter: router,
		Authenticator:       auth,
	}
}
