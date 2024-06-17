package main

import (
	"embed"

	"github.com/RicochetStudios/registry/client"
	"github.com/RicochetStudios/registry/client/wrapper"
)

var gameWrapper = wrapper.SimpleWrapper{
	StartProcess: wrapper.Process{
		Shell:  "/bin/bash",
		Script: "/start",
	},
	HealthProcess: wrapper.Process{
		Shell:  "/bin/bash",
		Script: "/health.sh",
	},
	ReadyCondition: wrapper.ReadyCondition{
		Message: `For help, type "help"`,
		Count:   1,
	},
	Timeouts: wrapper.ServerTimeouts{
		ReadyTimeout:    60,
		StopSoftTimeout: 10,
		StopHardTimeout: 15,
		HealthTimeout:   5,
	},
}

// embed the settings file into the binary.
// this means that the files are available as part of the package when imported.
//
//go:embed settings.yaml
var c embed.FS

func main() {
	// Configure settings.
	client.ConfigureSettings(c)

	// Check health flag.
	client.CheckHealth(&gameWrapper)

	// Start the main server lifecycle.
	client.Run(&gameWrapper)
}
