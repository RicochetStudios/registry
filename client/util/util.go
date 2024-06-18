package util

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Interceptor struct {
	Forward   io.Writer
	Intercept func(p []byte)
}

// Write will intercept the incoming stream, and forward
// the contents to its `forward` Writer.
func (i *Interceptor) Write(p []byte) (n int, err error) {
	if i.Intercept != nil {
		i.Intercept(p)
	}

	return i.Forward.Write(p)
}

// newSignalContext creates a new context that is cancelled when an interrupt or term signal is received.
func NewSignalContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigChan := make(chan os.Signal, 2)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		sig := <-sigChan
		InfoMessage("Received signal: '%s'\n" + sig.String())
		cancel()
	}()
	return ctx, cancel
}

// newReadyInterceptor creates a new Interceptor that will check for a ready message in the output.
// It will send a message to the ready channel when the ready message is found.
func NewReadyInterceptor(m string, c int, r chan bool) *Interceptor {
	// The number of ready statements required.
	var counter int = 0

	return &Interceptor{
		Forward: os.Stdout,
		Intercept: func(p []byte) {
			if counter >= c {
				return
			}

			str := strings.TrimSpace(string(p))
			// Checks if the ready message is in the output.
			if i := strings.Count(str, m); i > 0 {
				counter += i
				InfoMessage("Found ready statement: %d \n" + fmt.Sprint(counter))

				if counter >= c {
					InfoMessage("Moving to READY: %s \n" + str)
					r <- true
				}
			}
		}}
}
