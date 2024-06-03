package main

import (
	"github.com/RicochetStudios/registry/service"
	"github.com/RicochetStudios/registry/service/client"
)

var gameWrapper = client.SimpleWrapper{
	StartProcess: client.Process{
		Shell:  "/bin/bash",
		Script: "/start",
	},
	HealthProcess: client.Process{
		Shell:  "/bin/bash",
		Script: "/health.sh",
	},
	ReadyCondition: client.ReadyCondition{
		Message: `For help, type "help"`,
		Count:   1,
	},
	Timeouts: client.ServerTimeouts{
		ReadyTimeout:    60,
		StopSoftTimeout: 10,
		StopHardTimeout: 15,
		HealthTimeout:   5,
	},
}

func main() {
	// Check health flag.
	service.CheckHealth(&gameWrapper)

	// Start the main server lifecycle.
	service.Run(&gameWrapper)
}
