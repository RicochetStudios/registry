package service

import (
	"context"
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
