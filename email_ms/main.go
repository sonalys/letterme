package main

import (
	"context"

	"github.com/sonalys/letterme/domain/utils"
	"github.com/sonalys/letterme/email_ms/smtp"
)

func main() {
	// Context with cancel so we can stop all children from their inner loops after os.Interrupt.
	ctx, cancel := context.WithCancel(context.Background())

	smtp, err := smtp.InitServerFromEnv(ctx)
	if err != nil {
		panic(err)
	}

	go smtp.Listen()
	<-utils.GracefulShutdown()
	cancel()
	smtp.Shutdown()
}
