package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/sonalys/letterme/account_manager/controller"
	"github.com/sonalys/letterme/account_manager/persistence"
	"github.com/sonalys/letterme/account_manager/utils"
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
	utils.LoadFromEnv(persistence.MONGO_ENV, mongoConfig)
	mongo, err := persistence.NewMongo(ctx, mongoConfig)
	if err != nil {
		panic(err)
	}

	return &controller.Dependencies{
		Persistence: mongo,
	}
}
