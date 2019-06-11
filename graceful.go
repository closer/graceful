package graceful

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
)

var Trapped = errors.New("context trapped")

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
