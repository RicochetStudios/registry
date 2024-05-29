package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
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
func newSignalContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigChan := make(chan os.Signal, 2)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		sig := <-sigChan
		fmt.Printf("Received signal: '%s'\n", sig.String())
		cancel()
	}()
	return ctx, cancel
}
