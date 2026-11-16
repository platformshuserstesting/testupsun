// Command fanin-fanout is a small demo of the fan-out / fan-in pattern with
// goroutines: one producer sends jobs to several workers (fan-out) and their
// results are merged back into a single channel (fan-in).
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/platformshuserstesting/testupsun/fanio"
)

// maxWorkers caps -workers so a mistyped value is rejected instead of
// exhausting memory or overflowing.
const maxWorkers = 1000

func main() {
	workers := flag.Int("workers", 4, "number of worker goroutines (fan-out)")
	jobs := flag.Int("jobs", 20, "number of jobs to generate")
	delay := flag.Duration("delay", 50*time.Millisecond, "simulated work time per job")
	flag.Parse()

	if *workers < 1 || *workers > maxWorkers {
		fmt.Fprintf(os.Stderr, "error: -workers must be between 1 and %d\n", maxWorkers)
		os.Exit(2)
	}
	if *jobs < 0 {
		fmt.Fprintln(os.Stderr, "error: -jobs must not be negative")
		os.Exit(2)
	}
	if *delay < 0 {
		fmt.Fprintln(os.Stderr, "error: -delay must not be negative")
		os.Exit(2)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	nums := make([]int, *jobs)
	for i := range nums {
		nums[i] = i + 1
	}

	fmt.Printf("fan-out: %d jobs across %d workers (delay %s)\n", *jobs, *workers, *delay)

	in := fanio.Generate(ctx, nums...)
	outs := fanio.FanOut(ctx, in, *workers, *delay)
	results := fanio.Merge(ctx, outs...)

	perWorker := make([]int, *workers+1)
	total, sum := 0, 0
	for r := range results {
		fmt.Printf("worker %d processed %d -> %d\n", r.WorkerID, r.Input, r.Output)
		perWorker[r.WorkerID]++
		total++
		sum += r.Output
	}

	// A run is interrupted only if results are missing; a signal that arrives
	// after every result was received does not count.
	interrupted := total < *jobs
	if interrupted {
		fmt.Println("interrupted")
	}

	fmt.Printf("fan-in: received %d results, sum of outputs = %d\n", total, sum)
	for id := 1; id <= *workers; id++ {
		fmt.Printf("  worker %d handled %d jobs\n", id, perWorker[id])
	}

	if interrupted {
		stop()
		os.Exit(130)
	}
}
