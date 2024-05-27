package service

import (
	"context"
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
func newSignalContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := make(chan os.Signal, 2)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		cancel()
	}()
	return ctx
}
