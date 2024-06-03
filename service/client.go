package service

import (
	"context"
	"time"
)

// ServerWrapper is an interface that provides common operations with the server.
// It allows each game to have its own implementation of the server wrapper,
// with different logic relating to each specific game.
type ServerWrapper interface {
	// Start starts the server in the background.
	Start(context.Context) error
	// Wait returns once the server is ready.
	Wait() error
	// Wait returns once the server is ready.
	Serve(context.Context) error
	// Stop will attempt to shutdown the server gracefully.
	// Failing that, it will forcefully terminate.
	Stop()
	// Status returns the status of the server.
	Status() (string, error)
	// Healthy returns if the server is healthy.
	Healthy() (bool, error)
	// Logs returns the logs of the server.
	Logs() (string, error)
	// Backup invokes a backup of the server.
	Backup() error
}

// Process is a struct that contains the shell and script
// used to run a command.
// Can be used to start as server or check health of a server.
type Process struct {
	Shell  string
	Script string
}

// ReadyCondition is a struct that contains the message to check
// for in the output, and the number of times it should appear,
// before the server is considered ready.
type ReadyCondition struct {
	Message string
	Count   int
}

// ServerTimeouts is a struct that contains the timeouts for the server.
type ServerTimeouts struct {
	ReadyTimeout    time.Duration // Time to wait for the server to be ready when starting.
	StopSoftTimeout time.Duration // Time to wait for the server to stop gracefully.
	StopHardTimeout time.Duration // Time to wait for the server to stop forcefully.
	HealthTimeout   time.Duration // Time to wait for the server to check health.
}
