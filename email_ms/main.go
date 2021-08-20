package main

import (
	"context"

	"github.com/sonalys/letterme/domain/utils"
)

func main() {
	// Context with cancel so we can stop all children from their inner loops after os.Interrupt.
	_, cancel := context.WithCancel(context.Background())

	utils.GracefulShutdown(cancel)
}
