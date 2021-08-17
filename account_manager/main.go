package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/sonalys/letterme/account_manager/controller"
	"github.com/sonalys/letterme/account_manager/persistence"
)

func main() {
	// Context with cancel so we can stop all children from their inner loops after os.Interrupt.
	ctx, cancel := context.WithCancel(context.Background())

	// INFO: stop channel is needed for graceful shutdown implementation.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	mongo, err := persistence.NewMongo(ctx, &persistence.Configuration{})
	if err != nil {
		panic(err)
	}

	_, err = controller.NewService(ctx, &controller.Dependencies{
		Persistence: mongo,
	})
	if err != nil {
		panic(err)
	}

	<-stop
	cancel()
}
