package main

import (
	"context"
	"fmt"
	"log"

	sdk "agones.dev/agones/sdks/go"

	"github.com/RicochetStudios/registry/service/client"
	"github.com/RicochetStudios/registry/service/health"
)

func main() {
	// TODO: Allow extra arguments here.

	ctx := context.Background()

	// Connect to agones SDK.
	fmt.Println("Connecting to Agones with the SDK")
	s, err := sdk.NewSDK()
	if err != nil {
		log.Fatalf("Could not connect to Agones sdk: %v", err)
	}

	// Initialize the wrapper.
	fmt.Println("Initializing the server wrapper")
	wrapper := client.Sample{}
	fmt.Printf("Initialized the %v server wrapper", wrapper.GameName)

	// Start health check.
	fmt.Println("Starting health check")
	go health.DoHealth(&wrapper, s, ctx)

	// Start the server.
	fmt.Println("Starting the server")
	err = wrapper.Start()
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
	// Stop the server when the program ends.
	defer wrapper.Stop()

	// Wait for the server to be ready.
	fmt.Println("Waiting for the server to be ready")
	if err = wrapper.Wait(); err != nil {
		log.Fatalf("Server failed to reach ready state: %v", err)
	}

	// Serve to clients.
	fmt.Println("Serving the server")
	if err = wrapper.Serve(ctx); err != nil {
		log.Fatalf("Server failed to serve: %v", err)
	}
}
