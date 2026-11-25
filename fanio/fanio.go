// Package fanio demonstrates the fan-out / fan-in concurrency pattern with
// goroutines and channels.
//
// The pipeline has three stages:
//
//   - Generate: a single producer goroutine sends jobs on one channel.
//   - Worker / FanOut: several worker goroutines all read from that same
//     channel (fan-out), so the jobs are spread across them.
//   - Merge: the workers' result channels are combined into a single
//     channel (fan-in).
//
// Every stage watches ctx.Done(), so cancelling the context stops the whole
// pipeline without leaking goroutines.
package fanio

import (
	"context"
	"sync"
	"time"
)

// Result is the outcome of processing one input.
type Result struct {
	WorkerID int   // which worker handled the input
	Input    int   // the value read from the job channel
	Output   int64 // the computed value (Input squared)
}

// Generate is the producer stage. It sends each of nums on the returned
// channel and closes it when done or when ctx is cancelled.
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

// Worker reads jobs from in, squares each one and sends a Result on the
// returned channel. delay simulates work per item (zero means no delay).
// Several workers reading from the same in channel is the fan-out.
// The returned channel is closed when in is drained or ctx is cancelled.
func Worker(ctx context.Context, id int, in <-chan int, delay time.Duration) <-chan Result {
	out := make(chan Result)
	go func() {
		defer close(out)
		for {
			var n int
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				n = v
			case <-ctx.Done():
				return
			}
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
			case out <- Result{WorkerID: id, Input: n, Output: int64(n) * int64(n)}:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// FanOut starts n workers (with IDs 1..n) that all read from in, and returns
// their result channels.
func FanOut(ctx context.Context, in <-chan int, n int, delay time.Duration) []<-chan Result {
	chans := make([]<-chan Result, 0, n)
	for i := 1; i <= n; i++ {
		chans = append(chans, Worker(ctx, i, in, delay))
	}
	return chans
}

// Merge is the fan-in stage. It forwards values from every channel in chans
// onto a single returned channel, which is closed once all inputs are
// drained (or ctx is cancelled).
func Merge(ctx context.Context, chans ...<-chan Result) <-chan Result {
	out := make(chan Result)
	var wg sync.WaitGroup
	wg.Add(len(chans))
	for _, c := range chans {
		go func(c <-chan Result) {
			defer wg.Done()
			for {
				var r Result
				select {
				case v, ok := <-c:
					if !ok {
						return
					}
					r = v
				case <-ctx.Done():
					return
				}
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
