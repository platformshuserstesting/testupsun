// Package fanio demonstrates the fan-out / fan-in concurrency pattern with
// goroutines and channels.
//
// A single producer (Generate) emits jobs on a channel. Several workers
// (FanOut) read from that same channel concurrently, so the jobs are
// distributed among them: this is fan-out. Each worker writes its results to
// its own channel, and Merge combines all of those channels back into a
// single one: this is fan-in.
//
//	              +-> Worker 1 --+
//	Generate --> -+-> Worker 2 --+-> Merge --> consumer
//	              +-> Worker N --+
//
// Every stage watches ctx.Done() when sending, so cancelling the context stops
// the whole pipeline without leaking goroutines.
package fanio

import (
	"context"
	"sync"
	"time"
)

// Result is the output of a single job processed by a worker.
type Result struct {
	WorkerID int
	Input    int
	Output   int
}

// Generate is the producer stage. It sends each of nums on the returned
// unbuffered channel and closes it when all values have been sent or ctx is
// cancelled.
func Generate(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Worker is a single fan-out stage. It reads jobs from the shared in channel,
// "works" on each one (squares it after sleeping for delay), and sends a
// Result on the returned channel. The channel is closed when in is drained
// or ctx is cancelled.
func Worker(ctx context.Context, id int, in <-chan int, delay time.Duration) <-chan Result {
	out := make(chan Result)
	go func() {
		defer close(out)
		for n := range in {
			if delay > 0 {
				t := time.NewTimer(delay)
				select {
				case <-t.C:
				case <-ctx.Done():
					t.Stop()
					return
				}
			}
			select {
			case out <- Result{WorkerID: id, Input: n, Output: n * n}:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// FanOut starts n workers (with IDs 1..n) that all read from the same in
// channel and returns their output channels. Because the workers compete for
// values on in, each job is processed by exactly one worker.
func FanOut(ctx context.Context, in <-chan int, n int, delay time.Duration) []<-chan Result {
	chans := make([]<-chan Result, 0, n)
	for i := 1; i <= n; i++ {
		chans = append(chans, Worker(ctx, i, in, delay))
	}
	return chans
}

// Merge is the fan-in stage. It starts one goroutine per input channel that
// forwards values to a single output channel. A sync.WaitGroup tracks those
// goroutines, and the output channel is closed only once all of them have
// finished, i.e. when every input channel is closed or ctx is cancelled.
func Merge(ctx context.Context, chans ...<-chan Result) <-chan Result {
	out := make(chan Result)
	var wg sync.WaitGroup
	wg.Add(len(chans))
	for _, c := range chans {
		go func(c <-chan Result) {
			defer wg.Done()
			for r := range c {
				select {
				case out <- r:
				case <-ctx.Done():
					return
				}
			}
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
