package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	sdk "agones.dev/agones/sdks/go"
)

func Run(wrapper ServerWrapper) {
	// TODO: Allow extra arguments here.

	// Create a context that is cancelled when an interrupt or term signal is received.
	ctx := newSignalContext()

	// Connect to agones SDK.
	fmt.Println("Connecting to Agones with the SDK")
	s, err := sdk.NewSDK()
	if err != nil {
		log.Fatalf("Could not connect to Agones sdk: %v", err)
	}

	// Start health check.
	fmt.Println("Starting health check")
	go DoHealth(wrapper, s, ctx)

	// Start the server.
	fmt.Println("Starting the server")
	err = wrapper.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
	// Stop the server when the program ends.
	defer wrapper.Stop()

	// If an interupt or term signal is received, stop the server.

	// Wait for the server to be ready.
	fmt.Println("Waiting for the server to be ready")
	if err = wrapper.Wait(); err != nil {
		log.Fatalf("Server failed to reach ready state: %v", err)
	}
	s.Ready()

	// Serve to clients.
	fmt.Println("Serving the server")
	if err = wrapper.Serve(ctx); err != nil {
		log.Fatalf("Server failed to serve: %v", err)
	}
}

// newSignalContext creates a new context that is cancelled when an interrupt or term signal is received.
func newSignalContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		cancel()
	}()
	return ctx
}
