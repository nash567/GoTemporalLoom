package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

const gracefulShutDownTimeout = 10 * time.Second

func main() {
	// Set up a context that will be canceled when an os.Interrupt signal is received (e.g., Ctrl+C)
	// signal.NotifyContext creates a context that will be canceled when an os.Interrupt signal is caught
	// This also sets up listeners for the os.Interrupt signal, which need to be cleaned up later
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	// Ensure that the stop function is called when the main function exits,
	// which will clean up the signal notification setup and restore default signal handling behavior
	defer stop()

	// Initialize and start the origination application with the signal-aware context
	//originationApp := app.NewOriginationApp()
	//originationApp.Start(ctx)

	// Block the main function until the context is done, which means an interrupt signal was received
	<-ctx.Done()

	// Log that a shutdown signal has been received
	slog.Info("Received shutdown signal, starting graceful shutdown")

	// Reinstate the default signal handling behavior by calling stop
	// This stops listening for os.Interrupt signals in the custom way set up by signal.NotifyContext
	stop()

	// Set up a new context with a timeout for the graceful shutdown
	timeoutCtx, cancel := context.WithTimeout(context.Background(), gracefulShutDownTimeout)
	// Ensure that the cancel function is called when the graceful shutdown is complete
	defer cancel()

	// Create a goroutine to handle the timeout for the graceful shutdown
	go func() {
		<-timeoutCtx.Done()
		if err := timeoutCtx.Err(); errors.Is(err, context.DeadlineExceeded) {
			// If the graceful shutdown times out, log the error and forcefully exit the application
			slog.Error("Graceful shutdown timed out, shutting down forcefully", err)
			os.Exit(1)
		}
	}()

	// Begin the graceful shutdown of the origination application
	//originationApp.Stop(timeoutCtx)

}
