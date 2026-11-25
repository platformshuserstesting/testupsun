// Command testupsun demonstrates the fan-out / fan-in pattern with goroutines.
//
// A producer sends the numbers 1..n on one channel, several workers read from
// that channel concurrently (fan-out) and square each number, and their
// results are merged back into a single channel (fan-in).
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/platformshuserstesting/testupsun/fanio"
)

// maxN bounds -n so the sum of squares (about n³/3) stays well inside int64
// and the input slice stays small.
const maxN = 1_000_000

func main() {
	workers := flag.Int("workers", 3, "number of worker goroutines (fan-out)")
	n := flag.Int("n", 10, "process the numbers 1..n")
	delay := flag.Duration("delay", 0, "simulated work per item, e.g. 10ms")
	timeout := flag.Duration("timeout", 0, "stop the pipeline after this long (0 = no timeout)")
	flag.Parse()

	if *workers < 1 {
		fmt.Fprintln(os.Stderr, "error: -workers must be at least 1")
		os.Exit(1)
	}
	if *n < 0 || *n > maxN {
		fmt.Fprintf(os.Stderr, "error: -n must be between 0 and %d\n", maxN)
		os.Exit(1)
	}
	if *delay < 0 || *timeout < 0 {
		fmt.Fprintln(os.Stderr, "error: -delay and -timeout must not be negative")
		os.Exit(1)
	}

	ctx := context.Background()
	var cancel context.CancelFunc
	if *timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, *timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	nums := make([]int, *n)
	for i := range nums {
		nums[i] = i + 1
	}

	fmt.Printf("fan-out to %d workers, processing 1..%d\n", *workers, *n)

	jobs := fanio.Generate(ctx, nums...)
	results := fanio.Merge(ctx, fanio.FanOut(ctx, jobs, *workers, *delay)...)

	count, total := 0, 0
	for r := range results {
		fmt.Printf("worker %d: %d -> %d\n", r.WorkerID, r.Input, r.Output)
		count++
		total += r.Output
	}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) && count < *n {
		fmt.Printf("timeout after %s: received %d of %d results, partial sum %d\n", *timeout, count, *n, total)
		return
	}
	fmt.Printf("fan-in received %d results, sum of squares = %d\n", count, total)
}
