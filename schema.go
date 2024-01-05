package registry

import (
	"embed"
	"strings"

	"gopkg.in/yaml.v3"
)

type Sizes struct {
	XS Size `yaml:"xs"`
	S  Size `yaml:"s"`
	M  Size `yaml:"m"`
	L  Size `yaml:"l"`
	XL Size `yaml:"xl"`
}

type Size struct {
	Resources Resources `yaml:"resources"`
	Players   int       `yaml:"players"`
}

type Resources struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

type Network struct {
	Name     string `yaml:"name"`
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
}

type Setting struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Volume struct {
	Name  string `yaml:"name"`
	Path  string `yaml:"path"`
	Class string `yaml:"class"`
	Size  string `yaml:"size"`
}

type Probes struct {
	Command        []string `yaml:"command"`
	StartupProbe   Probe    `yaml:"startupProbe"`
	ReadynessProbe Probe    `yaml:"readynessProbe"`
	LivenessProbe  Probe    `yaml:"livenessProbe"`
}

type Probe struct {
	InitialDelaySeconds int `yaml:"initialDelaySeconds"`
	PeriodSeconds       int `yaml:"periodSeconds"`
	FailureThreshold    int `yaml:"failureThreshold"`
	SuccessThreshold    int `yaml:"successThreshold"`
	TimeoutSeconds      int `yaml:"timeoutSeconds"`
}

type Schema struct {
	Name     string          `yaml:"name"`
	Image    string          `yaml:"image"`
	URL      string          `yaml:"url"`
	Ratio    string          `yaml:"ratio"`
	Sizes    map[string]Size `yaml:"sizes"`
	Network  []Network       `yaml:"network"`
	Settings []Setting       `yaml:"settings"`
	Volumes  []Volume        `yaml:"volumes"`
	Probes   Probes          `yaml:"probes"`
}

// embed all schema files into the binary.
// this means that the files are available as part of the package when imported.
//
//go:embed schema/**.yaml
var c embed.FS

// GetSchema gets a game schema from a yaml file and stores it as a Schema.
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

// GetFileName corrects a file name if it does not already end with .yaml.
func getFileName(n string) string {
	s := strings.Split(n, ".")

	if s[len(s)-1] != ".yaml" {
		n = s[0] + ".yaml"
	}

	return n
}
