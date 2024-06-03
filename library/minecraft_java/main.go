package main

import (
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

func main() {
	// Check health flag.
	client.CheckHealth(&gameWrapper)

	// Start the main server lifecycle.
	client.Run(&gameWrapper)
}
