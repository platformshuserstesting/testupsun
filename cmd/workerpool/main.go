// Command workerpool demonstrates the workerpool package by squaring a range
// of numbers on a fixed number of worker goroutines.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/platformshuserstesting/testupsun/workerpool"
)

func main() {
	workers := flag.Int("workers", 4, "number of worker goroutines")
	jobs := flag.Int("jobs", 20, "number of jobs to submit")
	queue := flag.Int("queue", 10, "size of the task queue")
	flag.Parse()

	if err := run(*workers, *jobs, *queue); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(workers, jobs, queue int) error {
	if jobs < 0 {
		return fmt.Errorf("jobs must not be negative, got %d", jobs)
	}
	pool, err := workerpool.New(workers, queue)
	if err != nil {
		return err
	}

	start := time.Now()
	var mu sync.Mutex
	results := make([]int, jobs)

	ctx := context.Background()
	for i := 0; i < jobs; i++ {
		if err := pool.Submit(ctx, func() error {
			time.Sleep(50 * time.Millisecond) // simulate work
			mu.Lock()
			results[i] = i * i
			mu.Unlock()
			return nil
		}); err != nil {
			_ = pool.Close()
			return err
		}
	}

	if err := pool.Close(); err != nil {
		return err
	}
	elapsed := time.Since(start).Round(time.Millisecond)

	for i, r := range results {
		fmt.Printf("job %2d: %d^2 = %d\n", i, i, r)
	}
	fmt.Printf("processed %d jobs with %d workers in %s\n", jobs, workers, elapsed)
	return nil
}
