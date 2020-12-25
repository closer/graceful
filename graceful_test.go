package graceful_test

import (
	"context"
	"errors"
	"syscall"
	"testing"

	"github.com/closer/graceful"
)

func TestGracefulWithTrapping(t *testing.T) {
	ctx := graceful.WithTrap(context.Background(), syscall.SIGINT)

	// send signal
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	<-ctx.Done()
	err := ctx.Err()

	if !errors.Is(err, graceful.Trapped) {
		t.Errorf("expected graceful.Trapped but given %v", err)
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled but given %v", err)
	}
}

func TestGracefulWithCanceledChildContext(t *testing.T) {
	ctx := graceful.WithTrap(context.Background(), syscall.SIGINT)
	ctx, cancel := context.WithCancel(ctx)

	// cancel child context
	cancel()

	<-ctx.Done()
	err := ctx.Err()

	if errors.Is(err, graceful.Trapped) {
		t.Errorf("expected graceful.Trapped but given %v", err)
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled but given %v", err)
	}
}

func TestGracefulWithCanceledParentContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = graceful.WithTrap(ctx, syscall.SIGINT)

	// cancel parant context
	cancel()

	<-ctx.Done()
	err := ctx.Err()

	if errors.Is(err, graceful.Trapped) {
		t.Errorf("expected graceful.Trapped but given %v", err)
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled but given %v", err)
	}
}
