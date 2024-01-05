package registry

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Set desired result.
var want Schema = Schema{
	Name:  "minecraft_java",
	Image: "itzg/minecraft-server:latest",
	URL:   "https://github.com/itzg/docker-minecraft-server",
	Ratio: "1-2",
	Sizes: map[string]Size{
		"xs": {
			Resources: Resources{
				CPU:    "1000m",
				Memory: "2000Mi",
			},
			Players: 8,
		},
		"s": {
			Resources: Resources{
				CPU:    "1500m",
				Memory: "4000Mi",
			},
			Players: 16,
		},
		"m": {
			Resources: Resources{
				CPU:    "2000m",
				Memory: "8000Mi",
			},
			Players: 32,
		},
		"l": {
			Resources: Resources{
				CPU:    "3000m",
				Memory: "16000Mi",
			},
			Players: 64,
		},
		"xl": {
			Resources: Resources{
				CPU:    "4000m",
				Memory: "32000Mi",
			},
			Players: 128,
		},
	},
	Network: []Network{
		{
			Name:     "game",
			Port:     25565,
			Protocol: "tcp",
		},
	},
	Settings: []Setting{
		{
			Name:  "EULA",
			Value: "TRUE",
		},
		{
			Name:  "TYPE",
			Value: "{{ .modLoader }}",
		},
		{
			Name:  "MAX_PLAYERS",
			Value: "{{ .players }}",
		},
		{
			Name:  "MOTD",
			Value: "{{ .name }}",
		},
	},
	Volumes: []Volume{
		{
			Name:  "data",
			Path:  "/data",
			Class: "classic",
			Size:  "10Gi",
		},
	},
	Probes: Probes{
		Command: []string{"mc-health"},
		StartupProbe: Probe{
			FailureThreshold: 30,
			PeriodSeconds:    10,
		},
		ReadynessProbe: Probe{
			InitialDelaySeconds: 30,
			PeriodSeconds:       5,
			FailureThreshold:    20,
			SuccessThreshold:    3,
			TimeoutSeconds:      1,
		},
		LivenessProbe: Probe{
			InitialDelaySeconds: 30,
			PeriodSeconds:       5,
			FailureThreshold:    20,
			SuccessThreshold:    3,
			TimeoutSeconds:      1,
		},
	},
}

// TestGetSchema tests that GetSchema() will return the correct schema.
func TestGetSchema(t *testing.T) {
	// Call the function to test.
	got, err := GetSchema("minecraft_java.yaml")

	// Error if results are incorrect.
	if err != nil {
		t.Fatalf("GetSchema() returned an error: \n%v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("GetSchema() mismatch (-want +got):\n%s", diff)
	}
}

// TestGetSchemaWithoutFileExtension tests that GetSchema() will add the .yaml to a filename if it is missing.
func TestGetSchemaWithoutFileExtension(t *testing.T) {
	// Call the function to test.
	got, err := GetSchema("minecraft_java")

	// Error if results are incorrect.
	if err != nil {
		t.Fatalf("GetSchema() returned an error: \n%v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("GetSchema() mismatch (-want +got):\n%s", diff)
	}
}
