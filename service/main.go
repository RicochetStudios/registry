package service

import (
	"flag"
	"fmt"
	"log"
	"os"

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
	defer s.Shutdown()

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
	fmt.Println("Stopping serve")
}

// CheckHealth looks for a health flag and checks the health of the server.
func CheckHealth(wrapper ServerWrapper) {
	// Define the health flag.
	health := flag.Bool("health", false, "Check the health of the server")
	flag.Parse()

	if *health {
		healthy, err := wrapper.Healthy()
		if err != nil {
			fmt.Printf("Health check failed: %v\n", err)
			os.Exit(1)
		}

		if healthy {
			fmt.Println("Server is healthy")
		} else {
			fmt.Println("Server is unhealthy")
		}
		os.Exit(0)
	}
}
