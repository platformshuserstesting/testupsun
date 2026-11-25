package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/platformshuserstesting/testupsun/fanio"
)

const maxN = 1_000_000

const maxWorkers = 10_000

func main() {
	workers := flag.Int("workers", 3, "number of worker goroutines (fan-out)")
	n := flag.Int("n", 10, "process the numbers 1..n")
	delay := flag.Duration("delay", 0, "simulated work per item, e.g. 10ms")
	timeout := flag.Duration("timeout", 0, "stop the pipeline after this long (0 = no timeout)")
	flag.Parse()

	if *workers < 1 || *workers > maxWorkers {
		fmt.Fprintf(os.Stderr, "error: -workers must be between 1 and %d\n", maxWorkers)
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

	count, total := 0, int64(0)
	for r := range results {
		fmt.Printf("worker %d: %d -> %d\n", r.WorkerID, r.Input, r.Output)
		count++
		total += r.Output
	}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) && count < *n {
		fmt.Printf("timeout: %d of %d results, sum %d\n", count, *n, total)
		return
	}
	fmt.Printf("fan-in received %d results, sum of squares = %d\n", count, total)
}
