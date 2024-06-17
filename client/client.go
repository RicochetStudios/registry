package client

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/RicochetStudios/registry/client/util"
	"github.com/RicochetStudios/registry/client/wrapper"
	"gopkg.in/yaml.v2"

	sdk "agones.dev/agones/sdks/go"
)

func Run(wrapper wrapper.ServerWrapper) {
	// TODO: Allow extra arguments here.

	// Create a context that is cancelled when an interrupt or term signal is received.
	ctx, cancel := util.NewSignalContext()
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
func CheckHealth(wrapper wrapper.ServerWrapper) {
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

// ConfigureSettings sets game server environment variables,
// given an embedding of the settings file.
func ConfigureSettings(c embed.FS) error {
	const n = "settings.yaml"

	// Load the settings file contents.
	f, err := c.ReadFile(n)
	if err != nil {
		return fmt.Errorf("failed to load settings file embedding: %s", err)
	}

	// Create empty Settings to be are target of unmarshalling.
	var raw []struct {
		Name  string `yaml:"name"`
		Value string `yaml:"value"`
	}

	// Unmarshal the YAML file into empty Schema.
	if err := yaml.Unmarshal(f, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal settings file: %s", err)
	}

	// Convert the settings into a map.
	settings := make(map[string]string)
	for _, setting := range raw {
		settings[setting.Name] = setting.Value
	}

	// Set the environment variables.
	return util.RemapEnv(settings)
}
