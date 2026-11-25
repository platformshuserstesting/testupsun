// Package fanio implements a fan-out / fan-in pipeline with goroutines.
package fanio

import (
	"context"
	"sync"
	"time"
)

// Result is the outcome of processing one input.
type Result struct {
	WorkerID int
	Input    int
	Output   int64
}

// Generate sends each of nums on the returned channel.
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

// Worker squares each job from in, waiting delay per item.
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

// FanOut starts n workers that share in.
func FanOut(ctx context.Context, in <-chan int, n int, delay time.Duration) []<-chan Result {
	chans := make([]<-chan Result, 0, n)
	for i := 1; i <= n; i++ {
		chans = append(chans, Worker(ctx, i, in, delay))
	}
	return chans
}

// Merge combines chans into a single channel.
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
