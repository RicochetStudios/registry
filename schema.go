package registry

import (
	"embed"
	"strings"

	"gopkg.in/yaml.v3"
)

// Sizes are the available capacities for the game server.
// This changes the resources given to the server and the player limit.
//
// Sizes available are `xs`, `s`, `m`, `l`, `xl`.
type Sizes struct {
	XS Size `yaml:"xs"`
	S  Size `yaml:"s"`
	M  Size `yaml:"m"`
	L  Size `yaml:"l"`
	XL Size `yaml:"xl"`
}

// Size defines the resources allocated and player count for a single game size.
type Size struct {
	Resources Resources `yaml:"resources"` // CPU and Memory limits for specific size.
	Players   int       `yaml:"players"`   // Maximum number of players that can be supported on the resources.
}

// Resources defines the CPU and Memory limits for a game and is allocated to a size.
//
// Cpu and memory specify the limits of the resources provided to the server.
// The server can run at any amount of resources below and up to this limit, but not above.
//
// Functions identically to and provided as a string of Kubernetes resource units, see:
// https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#resource-units-in-kubernetes.
type Resources struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

// Network defines a set of port and protocols to expose from the game server runs on.
// This makes the game server accessible from the method is used to connect to it.
//
// Some games can have multiple services to expose, to function correctly.
type Network struct {
	Name     string `yaml:"name"`     // The name of the network configuration.
	Port     int    `yaml:"port"`     // The port of the service.
	Protocol string `yaml:"protocol"` // The protocol of the service.
}

// Setting defines a static environment variable specific to the game server image.
//
// Set as a key value pair.
type Setting struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Volume is storage data for the game server.
// This is used to persist the game server data between restarts.
//
// Provided as a string of Kubernetes capacity units, see:
// https://kubernetes.io/docs/concepts/storage/persistent-volumes/#capacity
type Volume struct {
	Name string `yaml:"name"` // The name of the volume.
	Path string `yaml:"path"` // The path of the data.
	Size string `yaml:"size"` // The maximum size of storage.
}

// Probes are used to check the health of the game server.
type Probes struct {
	Command        []string `yaml:"command"`        // A list of commands to run to check the health of the game server.
	StartupProbe   Probe    `yaml:"startupProbe"`   // The startup probe checks if the game server has started correctly.
	ReadinessProbe Probe    `yaml:"readinessProbe"` // The readiness probe checks if the game server is ready to accept connections.
	LivenessProbe  Probe    `yaml:"livenessProbe"`  // The liveness probe checks if the game server is still running.
}

// Probe defines the configuration for a probe.
type Probe struct {
	InitialDelaySeconds int `yaml:"initialDelaySeconds"` // Seconds after the container has started before probes are initiated.
	PeriodSeconds       int `yaml:"periodSeconds"`       // Seconds between probe checks.
	FailureThreshold    int `yaml:"failureThreshold"`    // Consecutive failures of the probe before the container is restarted.
	SuccessThreshold    int `yaml:"successThreshold"`    // Consecutive successes of the probe for the contiainer to be deemed healthy.
	TimeoutSeconds      int `yaml:"timeoutSeconds"`      // Seconds after which the probe times out.
}

// Schema is the structure of a game servers settings.
type Schema struct {
	Name     string          `yaml:"name"`     // Name of the game.
	Image    string          `yaml:"image"`    // Container image of the game server.
	URL      string          `yaml:"url"`      // Metadata providing a link to find out more about the image.
	Ratio    string          `yaml:"ratio"`    // Ratio of CPU to Memory for the game server.
	Sizes    map[string]Size `yaml:"sizes"`    // Sizes available for the game server, their resources and player limits.
	Network  []Network       `yaml:"network"`  // Network configuration for the game server.
	Settings []Setting       `yaml:"settings"` // Environment variables for the game server container image.
	Volumes  []Volume        `yaml:"volumes"`  // Persistent storage volumes for the game server data.
	Probes   Probes          `yaml:"probes"`   // Health checks for the game server.
}

// embed all schema files into the binary.
// this means that the files are available as part of the package when imported.
//
//go:embed schema/**.yaml
var c embed.FS

// GetSchema gets a game Schema from a yaml file, when given the name of the game.
func GetSchema(n string) (Schema, error) {
	// Correct the file name if it does not already end with .yaml.
	n = getFileName(n)

	// Load the schema file contents.
	f, err := c.ReadFile("schema/" + n)
	if err != nil {
		return Schema{}, err
	}

	// Create an empty Schema to be are target of unmarshalling.
	var schema Schema

	// Unmarshal the YAML file into empty Schema.
	if err := yaml.Unmarshal(f, &schema); err != nil {
		return Schema{}, err
	}

	return schema, nil
}

// getFileName corrects a file name if it does not already end with .yaml.
func getFileName(n string) string {
	s := strings.Split(n, ".")

	if s[len(s)-1] != ".yaml" {
		n = s[0] + ".yaml"
	}

	return n
}
