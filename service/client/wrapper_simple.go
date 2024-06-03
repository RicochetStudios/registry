package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/RicochetStudios/registry/service/util"
)

// SimpleWrapper is a struct that implements the ServerWrapper interface.
// It contains the necessary information to start, stop, and interact with
// common, basic dedicated video game server.
type SimpleWrapper struct {
	Cmd    *exec.Cmd
	Ready  chan bool
	Ctx    context.Context
	Cancel context.CancelFunc
	Stdout bytes.Buffer
	Stderr bytes.Buffer

	StartProcess   Process
	HealthProcess  Process
	ReadyCondition ReadyCondition
	Timeouts       ServerTimeouts
}

// init initializes the server object.
func (m *SimpleWrapper) init(ctx context.Context) {
	// Create a context for the server.
	m.Ctx, m.Cancel = context.WithCancel(ctx)

	// Create a channel to check if the server is ready.
	m.Ready = make(chan bool, 1)
}

// Start starts the server in the background.
func (m *SimpleWrapper) Start(ctx context.Context) error {
	var (
		startShell   = m.StartProcess.Shell
		startScript  = m.StartProcess.Script
		readyMessage = m.ReadyCondition.Message
		readyCount   = m.ReadyCondition.Count
	)

	// Initialize the server object.
	m.init(ctx)

	// Assign the command to the server.
	m.Cmd = exec.Command(startShell, startScript)
	// If the context is cancelled, stop the server.
	go func() {
		<-m.Ctx.Done()
		m.Stop()
	}()

	// Intercept the stdout to check if the server is ready.
	m.Cmd.Stderr = io.MultiWriter(&m.Stderr, &util.Interceptor{Forward: os.Stderr})
	m.Cmd.Stdout = io.MultiWriter(&m.Stdout,
		util.NewReadyInterceptor(readyMessage, readyCount, m.Ready))

	// Start the server.
	err := m.Cmd.Start()
	if err != nil {
		return fmt.Errorf("the start script failed with error: %v\nExecuting the following script:\n%v", err, startScript)
	}

	return nil
}

// Wait returns once the server is ready.
func (m *SimpleWrapper) Wait() error {
	// Create a timeout for the server to be ready.
	var readyTimeout = m.Timeouts.ReadyTimeout

	// Return false if the timeout is exeeded.
	go func() {
		time.Sleep(readyTimeout * time.Second)
		m.Ready <- false
	}()

	// Wait for the server to be ready.
	rdy := <-m.Ready
	if rdy {
		return nil
	}

	return fmt.Errorf("server failed to reach ready state within timeout of %d seconds", int(m.Timeouts.ReadyTimeout))
}

// Serve serves the server to clients.
func (m *SimpleWrapper) Serve(context.Context) error {
	// Wait for the server to exit.
	if err := m.Cmd.Wait(); err != nil {
		return fmt.Errorf("server crashed: %v", err)
	}

	return nil
}

// Stop will attempt to shutdown the server gracefully.
// Failing that, it will forcefully terminate.
func (m *SimpleWrapper) Stop() {
	// Check if the server is already stopped.
	if m.Cmd.ProcessState != nil {
		return
	}
	fmt.Println("Stopping the server")

	// Define the timeouts for stopping the server.
	var (
		stopSoftTimeout = m.Timeouts.StopSoftTimeout
		stopHardTimeout = m.Timeouts.StopHardTimeout
	)

	// Release the server process after a timeout.
	// This is a last resort to prevent a zombie process.
	go func() {
		time.Sleep(stopHardTimeout * time.Second)
		if err := m.Cmd.Process.Release(); err != nil {
			fmt.Printf("failed to release server process: %v\n", err)
		}
	}()

	// Attempt to backup the server before stopping.
	if err := m.Backup(); err != nil {
		fmt.Printf("Failed to backup server: %v\n", err)
	} else {
		fmt.Println("Server backed up successfully")
	}

	// Attempt to stop the server gracefully
	// by sending a terminating signal.
	//
	// It's important to use a terminating signal,
	// as most applications recognize this signal and,
	// the exec cmd will exit correctly.
	if err := m.Cmd.Process.Signal(syscall.SIGTERM); err != nil {
		fmt.Printf("Failed to stop server gracefully: %v\n", err)
	}

	// Check if the server has stopped in the timeout.
	for i := 0; i < int(stopSoftTimeout); i++ {
		// Check if the server has stopped.
		// ProcessState.Exited() will only return true if
		// the process was killed or a SIGTERM was sent.
		if m.Cmd.ProcessState != nil && m.Cmd.ProcessState.Exited() {
			fmt.Printf("Server stopped gracefully after %d seconds\n", i)
			return
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("Server failed to stop gracefully within soft timeout of %d seconds\n", int(stopSoftTimeout))

	// If the server has not stopped, forcefully terminate it.
	if err := m.Cmd.Process.Kill(); err != nil {
		fmt.Printf("Failed to stop server forcefully: %v\n", err)
	} else {
		// Wait for the server to exit.
		m.Cmd.Wait()
		fmt.Println("Server stopped forcefully")
		return
	}

	fmt.Printf("Releasing server resources\n")
}

// Status returns the status of the server.
func (m *SimpleWrapper) Status() (string, error) {
	return "OK", nil
}

// Healthy returns if the server is healthy.
func (m *SimpleWrapper) Healthy() (bool, error) {
	// Define the health check variables.
	var (
		healthTimeout = m.Timeouts.HealthTimeout
		healthShell   = m.HealthProcess.Shell
		healthScript  = m.HealthProcess.Script
	)

	// Create a context to cancel the health check after a timeout.
	ctx, cancel := context.WithTimeout(m.Ctx, healthTimeout*time.Second)
	defer cancel()

	// Check if the server is healthy,
	// by sending a command to a binary.
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
func (m *SimpleWrapper) Logs() (string, error) {
	return "Logs", nil
}

// Backup invokes a backup of the server.
func (m *SimpleWrapper) Backup() error {
	return nil
}
