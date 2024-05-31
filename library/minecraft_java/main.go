package main

import (
	"github.com/RicochetStudios/registry/service"
)

var GameWrapper = service.BasicWrapper{
	StartProccess: service.StartProccess{
		Shell:  "/bin/bash",
		Script: "/start",
	},
	HealthProccess: service.HealthProccess{
		Shell:  "/bin/bash",
		Script: "/health.sh",
	},
	ReadyCondition: service.ReadyCondition{
		Message: `For help, type "help"`,
		Count:   1,
	},
	Timeouts: service.ServerTimeouts{
		ReadyTimeout:    60,
		StopSoftTimeout: 10,
		StopHardTimeout: 15,
		HealthTimeout:   5,
	},
}

func main() {
	// Check health flag.
	service.CheckHealth(&GameWrapper)

	// Start the main server lifecycle.
	service.Run(&GameWrapper)
}
