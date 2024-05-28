package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/RicochetStudios/registry/service"
)

// MinecraftJava is a struct that implements the ServerWrapper interface.
// It contains the necessary information to start, stop, and interact with a Minecraft Java server.
type MinecraftJava struct {
	Cmd    *exec.Cmd
	Ready  bool
	Ctx    context.Context
	Cancel context.CancelFunc
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

var GameWrapper MinecraftJava = MinecraftJava{}

func main() {
	// Define the health flag.
	health := flag.Bool("health", false, "Check the health of the server")
	flag.Parse()

	if *health {
		healthy, err := GameWrapper.Healthy()
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

	// Start the main server lifecycle.
	service.Run(&GameWrapper)
}

// Start starts the server in the background.
func (m *MinecraftJava) Start(ctx context.Context) error {
	const (
		startShell  = "/bin/bash"
		startScript = "/start"
	)

	// Create a context for the server.
	m.Ctx, m.Cancel = context.WithCancel(ctx)

	// Assign the command to the server.
	m.Cmd = exec.CommandContext(m.Ctx, startShell, startScript)

	m.Cmd.Cancel = func() error {
		m.Stop()
		if err := m.Cmd.Process.Release(); err != nil {
			return fmt.Errorf("failed to release server process: %v", err)
		}
		return nil
	}

	// The number of ready statements required.
	var readyCount int = 0

	// Intercept the stdout to check if the server is ready.
	m.Cmd.Stderr = io.MultiWriter(&m.Stderr, &service.Interceptor{Forward: os.Stderr})
	m.Cmd.Stdout = io.MultiWriter(&m.Stdout, &service.Interceptor{
		Forward: os.Stdout,
		Intercept: func(p []byte) {
			if readyCount >= 1 {
				return
			}

			str := strings.TrimSpace(string(p))
			// Minecraft Java will say "[Server] Startup Done" once ready.
			if count := strings.Count(str, `For help, type "help"`); count > 0 {
				readyCount += count
				fmt.Printf("Found ready statement: %d \n", readyCount)

				if readyCount <= 1 {
					fmt.Printf("Moving to READY: %s \n", str)
					m.Ready = true
				}
			}
		}})

	// Start the server.
	err := m.Cmd.Start()
	if err != nil {
		return fmt.Errorf("the start script failed with error: %v\nExecuting the following script:\n%v", err, startScript)
	}

	return nil
}

// Wait returns once the server is ready.
func (m *MinecraftJava) Wait() error {
	// The maximum time to wait for the server to be ready.
	// TODO: Review this value.
	const readyTimeout = 120

	// Wait for the server to be ready.
	for i := 0; i < readyTimeout; i++ {
		if m.Ready {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	// TODO: Add a more descriptive error message, with stdout and stderr.
	return fmt.Errorf("server failed to reach ready state within timeout of %d seconds", readyTimeout)
}

// Serve serves the server to clients.
func (m *MinecraftJava) Serve(context.Context) error {
	// Wait for the server to exit.
	if err := m.Cmd.Wait(); err != nil {
		return fmt.Errorf("server crashed: %v", err)
	}

	return nil
}

// Stop will attempt to shutdown the server gracefully.
// Failing that, it will forcefully terminate.
func (m *MinecraftJava) Stop() {
	fmt.Println("Stopping the server")
	const stopSoftTimeout = 10
	const stopHardTimeout = stopSoftTimeout + 5

	// Cancel the context to release the server process.
	// This is a last resort.
	go func() {
		time.Sleep(stopHardTimeout * time.Second)
		m.Cancel()
	}()

	// Attempt to backup the server before stopping.
	if err := m.Backup(); err != nil {
		fmt.Printf("Failed to backup server: %v\n", err)
	} else {
		fmt.Println("Server backed up successfully")
	}

	// Attempt to stop the server gracefully
	// by sending an interrupt signal.
	if err := m.Cmd.Process.Signal(os.Interrupt); err != nil {
		fmt.Printf("Failed to stop server gracefully: %v\n", err)
	}

	// Check if the server has stopped in the timeout.
	for i := 0; i < stopSoftTimeout; i++ {
		if m.Cmd.ProcessState != nil && m.Cmd.ProcessState.Exited() {
			fmt.Println("Server stopped gracefully")
			return
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("Server failed to stop gracefully within soft timeout of %d seconds\n", stopSoftTimeout)

	// If the server has not stopped, forcefully terminate it.
	if err := m.Cmd.Process.Kill(); err != nil {
		fmt.Printf("Failed to stop server forcefully: %v\n", err)
	} else {
		// Wait for the server to exit.
		m.Cmd.Wait()
		fmt.Println("Server stopped forcefully")
		return
	}

	fmt.Printf("Server failed to stop forcefully, releasing resources\n")
}

// Status returns the status of the server.
func (m *MinecraftJava) Status() (string, error) {
	return "OK", nil
}

// Healthy returns if the server is healthy.
func (m *MinecraftJava) Healthy() (bool, error) {
	const (
		healthTimeout = 5
		healthShell   = "/bin/bash"
		healthScript  = "/health.sh"
	)

	// Create a context to cancel the health check after a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), healthTimeout*time.Second)
	defer cancel()

	// Check if the server is healthy,
	// by sending a command to the mc-health binary.
	cmd := exec.CommandContext(ctx, healthShell, healthScript)

	// Run the command.
	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("health check failed to complete: %v", err)
	}

	// Check the output of the command.
	if cmd.ProcessState.ExitCode() == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// Logs returns the logs of the server.
func (m *MinecraftJava) Logs() (string, error) {
	return "Logs", nil
}

// Backup invokes a backup of the server.
func (m *MinecraftJava) Backup() error {
	return nil
}
