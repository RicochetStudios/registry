package service

import "io"

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
