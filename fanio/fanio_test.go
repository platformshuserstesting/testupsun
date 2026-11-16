package fanio

import (
	"context"
	"sort"
	"testing"
	"time"
)

func seq(n int) []int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i + 1
	}
	return nums
}

func TestPipelineProducesAllSquares(t *testing.T) {
	tests := []struct {
		name    string
		workers int
		jobs    int
	}{
		{"one worker", 1, 20},
		{"three workers", 3, 20},
		{"ten workers", 10, 20},
		{"more workers than jobs", 10, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			nums := seq(tt.jobs)
			results := Merge(ctx, FanOut(ctx, Generate(ctx, nums...), tt.workers, 0)...)

			var got []int
			for r := range results {
				if r.Output != r.Input*r.Input {
					t.Errorf("result %+v: output is not the square of input", r)
				}
				if r.WorkerID < 1 || r.WorkerID > tt.workers {
					t.Errorf("result %+v: worker ID out of range 1..%d", r, tt.workers)
				}
				got = append(got, r.Output)
			}

			if len(got) != len(nums) {
				t.Fatalf("got %d results, want %d", len(got), len(nums))
			}
			sort.Ints(got)
			for i, n := range nums {
				if got[i] != n*n {
					t.Fatalf("sorted results[%d] = %d, want %d", i, got[i], n*n)
				}
			}
		})
	}
}

func TestMergeNoChannelsClosesImmediately(t *testing.T) {
	out := Merge(context.Background())
	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("expected closed channel, got a value")
		}
	case <-time.After(time.Second):
		t.Fatal("Merge with no channels did not close its output")
	}
}

func TestCancelStopsPipeline(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	results := Merge(ctx, FanOut(ctx, Generate(ctx, seq(1000)...), 4, 10*time.Millisecond)...)

	// Read one result, then cancel and make sure the output closes promptly.
	<-results
	cancel()

	done := make(chan struct{})
	go func() {
		for range results {
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("pipeline did not stop after cancellation")
	}
}

func TestWorkIsSpreadAcrossWorkers(t *testing.T) {
	ctx := context.Background()
	results := Merge(ctx, FanOut(ctx, Generate(ctx, seq(20)...), 4, 5*time.Millisecond)...)

	seen := map[int]bool{}
	for r := range results {
		seen[r.WorkerID] = true
	}
	if len(seen) < 2 {
		t.Fatalf("expected jobs on more than one worker, got workers %v", seen)
	}
}
