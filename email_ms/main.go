package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sonalys/letterme/domain/messaging"
	"github.com/sonalys/letterme/domain/messaging/rabbitmq"
	"github.com/sonalys/letterme/domain/persistence/mongo"
	"github.com/sonalys/letterme/domain/utils"
	"github.com/sonalys/letterme/email_ms/controller"
	"github.com/sonalys/letterme/email_ms/handler"
	"github.com/sonalys/letterme/email_ms/smtp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	smtp, err := smtp.InitServerFromEnv(ctx)
	if err != nil {
		panic(err)
	}

	initialize(ctx, smtp)

	go smtp.Listen()
	<-utils.GracefulShutdown()
	cancel()
	smtp.Shutdown()
}

func initialize(ctx context.Context, smtp *smtp.Server) {
	mongoConfig := new(mongo.Configuration)
	if err := utils.LoadFromEnv(mongo.ConfigEnv, mongoConfig); err != nil {
		logrus.Panicf("failed to initialize mongoConfig from env: %s", err)
	}

	mongo, err := mongo.NewClient(ctx, mongoConfig)
	if err != nil {
		panic(err)
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

	router, err := messaging.NewRouter(ctx, routerConfig, &messaging.Dependencies{
		Messenger: rabbit,
	})
	if err != nil {
		panic(err)
	}

	controllerConfig := &controller.Configuration{}
	svc, err := controller.NewService(ctx, controllerConfig, &controller.Dependencies{
		EventRouter: router,
		Persistence: mongo,
		Messenger:   rabbit,
	})
	if err != nil {
		panic(err)
	}

	smtp.AddMiddlewares(svc.ValidateEmailMiddleware)
	handler.RegisterHandlers(router, svc)
}
