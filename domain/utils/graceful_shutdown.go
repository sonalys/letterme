package utils

import (
	"os"
	"os/signal"
)

// GracefulShutdown is used to stop the main go routine after os request.
func GracefulShutdown(cancelContext func()) {
	// INFO: stop channel is needed for graceful shutdown implementation.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	cancelContext()
}
