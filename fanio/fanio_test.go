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

func collect(t *testing.T, ch <-chan Result) []Result {
	t.Helper()
	var got []Result
	timeout := time.After(5 * time.Second)
	for {
		select {
		case r, ok := <-ch:
			if !ok {
				return got
			}
			got = append(got, r)
		case <-timeout:
			t.Fatal("timed out waiting for merged channel to close")
		}
	}
}

func TestEveryInputProcessedExactlyOnce(t *testing.T) {
	ctx := context.Background()
	nums := seq(100)

	results := collect(t, Merge(ctx, FanOut(ctx, Generate(ctx, nums...), 4, 0)...))

	if len(results) != len(nums) {
		t.Fatalf("got %d results, want %d", len(results), len(nums))
	}
	outputs := make([]int, len(results))
	for i, r := range results {
		if r.Output != r.Input*r.Input {
			t.Errorf("worker %d: %d -> %d, want %d", r.WorkerID, r.Input, r.Output, r.Input*r.Input)
		}
		outputs[i] = r.Output
	}
	sort.Ints(outputs)
	for i, n := range nums {
		if outputs[i] != n*n {
			t.Fatalf("sorted outputs[%d] = %d, want %d", i, outputs[i], n*n)
		}
	}
}

func TestWorkIsSpreadAcrossWorkers(t *testing.T) {
	ctx := context.Background()
	const workers = 3

	results := collect(t, Merge(ctx, FanOut(ctx, Generate(ctx, seq(30)...), workers, 5*time.Millisecond)...))

	ids := map[int]bool{}
	for _, r := range results {
		if r.WorkerID < 1 || r.WorkerID > workers {
			t.Fatalf("unexpected worker id %d", r.WorkerID)
		}
		ids[r.WorkerID] = true
	}
	if len(ids) < 2 {
		t.Fatalf("expected results from more than one worker, got worker ids %v", ids)
	}
}

func TestCancelClosesMergedChannel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	merged := Merge(ctx, FanOut(ctx, Generate(ctx, seq(1000)...), 3, 10*time.Millisecond)...)

	// Read a single result, then cancel and stop reading promptly.
	select {
	case <-merged:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for first result")
	}
	cancel()

	results := collect(t, merged)
	if len(results) >= 999 {
		t.Fatalf("expected cancellation to stop the pipeline early, got %d more results", len(results))
	}
}

func TestMergeWithNoChannelsClosesImmediately(t *testing.T) {
	select {
	case _, ok := <-Merge(context.Background()):
		if ok {
			t.Fatal("expected closed channel, got a value")
		}
	case <-time.After(time.Second):
		t.Fatal("Merge with no channels did not close")
	}
}
