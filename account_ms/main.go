package main

import (
	"context"

	"github.com/sonalys/letterme/account_ms/controller"
	"github.com/sonalys/letterme/domain/utils"
)

func main() {
	// Context with cancel so we can stop all children from their inner loops after os.Interrupt.
	ctx, cancel := context.WithCancel(context.Background())

	if _, err := controller.InitializeFromEnv(ctx); err != nil {
		panic(err)
	}

	<-utils.GracefulShutdown()
	cancel()
}
