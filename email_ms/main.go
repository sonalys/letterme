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
	// Context with cancel so we can stop all children from their inner loops after os.Interrupt.
	ctx, cancel := context.WithCancel(context.Background())

	initialize(ctx)

	smtp, err := smtp.InitServerFromEnv(ctx)
	if err != nil {
		panic(err)
	}

	go smtp.Listen()
	<-utils.GracefulShutdown()
	cancel()
	smtp.Shutdown()
}

// initialize starts all the ms dependencies and sub-routines,
// if it fails, the ms will panic.
func initialize(ctx context.Context) {
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
	}

	router, err := messaging.NewRouter(ctx, routerConfig, &messaging.Dependencies{
		Messaging: rabbit,
	})
	if err != nil {
		panic(err)
	}

	controllerConfig := &controller.Configuration{}
	svc, err := controller.NewService(ctx, controllerConfig, &controller.Dependencies{
		Router:      router,
		Persistence: mongo,
		Messaging:   rabbit,
	})
	if err != nil {
		panic(err)
	}

	handler.RegisterHandlers(router, svc)
}
