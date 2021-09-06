package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sonalys/letterme/account_ms/controller"
	"github.com/sonalys/letterme/account_ms/handler"
	"github.com/sonalys/letterme/domain/cryptography"
	"github.com/sonalys/letterme/domain/messaging"
	"github.com/sonalys/letterme/domain/messaging/rabbitmq"
	"github.com/sonalys/letterme/domain/persistence/mongo"
	"github.com/sonalys/letterme/domain/utils"
)

func main() {
	// Context with cancel so we can stop all children from their inner loops after os.Interrupt.
	ctx, cancel := context.WithCancel(context.Background())

	initialize(ctx)

	<-utils.GracefulShutdown()
	cancel()
}

func initialize(ctx context.Context) {
	mongoConfig := new(mongo.Configuration)
	if err := utils.LoadFromEnv(mongo.ConfigEnv, mongoConfig); err != nil {
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

	cryptoRouter, err := cryptography.NewCryptoRouter(cryptographicConfig)
	if err != nil {
		panic(err)
	}

	authConfig := new(cryptography.AuthConfiguration)
	if err := utils.LoadFromEnv(cryptography.JWT_AUTH_ENV, authConfig); err != nil {
		logrus.Panicf("failed to initialize authConfig from env: %s", err)
	}
	auth := cryptography.NewJWTAuthenticator(authConfig)

	svcConfig := new(controller.Configuration)
	if err := utils.LoadFromEnv(controller.ServiceConfigurationEnv, svcConfig); err != nil {
		logrus.Panicf("failed to initialize authConfig from env: %s", err)
	}

	rabbitConfig := new(rabbitmq.Configuration)
	if err := utils.LoadFromEnv(rabbitmq.ConfigEnv, mongoConfig); err != nil {
		logrus.Panicf("failed to initialize mongoConfig from env: %s", err)
	}

	rabbit, err := rabbitmq.NewClient(rabbitConfig)
	if err != nil {
		panic(err)
	}

	// TODO: fix the configuration here.
	routerConfig := &messaging.Configuration{
		ResponseTimeout: time.Second,
		ResponseChannel: "email_ms",
	}

	eventRouter, err := messaging.NewRouter(ctx, routerConfig, &messaging.Dependencies{
		Messaging: rabbit,
	})
	if err != nil {
		panic(err)
	}

	svc, err := controller.NewService(ctx, svcConfig, &controller.Dependencies{
		Persistence:         mongo,
		CryptographicRouter: cryptoRouter,
		Authenticator:       auth,
	})
	if err != nil {
		panic(err)
	}

	handler.RegisterHandlers(eventRouter, svc)
}
