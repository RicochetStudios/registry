package service

import (
	"fmt"
	"log"

	sdk "agones.dev/agones/sdks/go"
)

func Run(wrapper ServerWrapper) {
	// TODO: Allow extra arguments here.

	// Create a context that is cancelled when an interrupt or term signal is received.
	ctx, cancel := newSignalContext()
	// Stop the server when the program ends.
	defer cancel()

	// Connect to agones SDK.
	fmt.Println("Connecting to Agones with the SDK")
	s, err := sdk.NewSDK()
	if err != nil {
		log.Fatalf("Could not connect to Agones sdk: %v", err)
	}

	// Start the server.
	fmt.Println("Starting the server")
	err = wrapper.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
	// Stop the server when the program ends.
	// defer wrapper.Stop()

	// Wait for the server to be ready.
	fmt.Println("Waiting for the server to be ready")
	if err = wrapper.Wait(); err != nil {
		log.Fatalf("Server failed to reach ready state: %v", err)
	}
	s.Ready()

	// Start health check.
	fmt.Println("Starting health check")
	go DoHealth(wrapper, s, ctx)

	// Serve to clients.
	fmt.Println("Beginning to serve")
	if err = wrapper.Serve(ctx); err != nil {
		log.Fatalf("Server failed to serve: %v", err)
	}
}
