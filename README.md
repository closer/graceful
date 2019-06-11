# graceful

[![GoDoc](https://godoc.org/github.com/closer/graceful?status.svg)](https://godoc.org/github.com/closer/graceful)

Graceful is Go package for signal trapping based on context.

## Usage

```go
package main

import (
	"context"
	"log"
	"syscall"

	"github.com/closer/graceful"
)

func main() {
	log.Println("Starting up...")

	ctx := graceful.WithTrap(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	if err := Worker(ctx); err != nil {
		log.Fatal("Worker error", err)
	}

	log.Println("Shutting down...")
}
```
