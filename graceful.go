package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

// Trapped is the error returned by Context.Err when the context is trapped.
// It wraps context.Canceled.
var Trapped = &trapped{}

type trapped struct{}

func (t *trapped) Error() string {
	return "context trapped"
}

func (t *trapped) Unwrap() error {
	return context.Canceled
}

// WithTrap returns a copy of parent with a new Done channel. The returned
// context's Done channel is closed when the specified signal is trapped
// or when the parent context's Done channel is closed, whichever happens first.
func WithTrap(parent context.Context, sigs ...os.Signal) context.Context {
	trap := make(chan os.Signal, 1)
	signal.Notify(trap, sigs...)

	c := &trapCtx{Context: parent}

	go func() {
		select {
		case <-parent.Done():
			c.cancel(parent.Err())
		case <-trap:
			c.cancel(Trapped)
		}
	}()

	return c
}

var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

type trapCtx struct {
	context.Context

	mu   sync.Mutex
	done chan struct{}
	err  error
}

func (c *trapCtx) Done() <-chan struct{} {
	c.mu.Lock()
	if c.done == nil {
		c.done = make(chan struct{})
	}
	d := c.done
	c.mu.Unlock()
	return d
}

func (c *trapCtx) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

func (c *trapCtx) cancel(err error) {
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return
	}
	c.err = err
	if c.done == nil {
		c.done = closedchan
	} else {
		close(c.done)
	}
	c.mu.Unlock()
}
