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
	outputs := make([]int64, len(results))
	for i, r := range results {
		want := int64(r.Input) * int64(r.Input)
		if r.Output != want {
			t.Errorf("worker %d: %d -> %d, want %d", r.WorkerID, r.Input, r.Output, want)
		}
		outputs[i] = r.Output
	}
	sort.Slice(outputs, func(i, j int) bool { return outputs[i] < outputs[j] })
	for i, n := range nums {
		if outputs[i] != int64(n)*int64(n) {
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
	const workers = 3
	ctx, cancel := context.WithCancel(context.Background())
	merged := Merge(ctx, FanOut(ctx, Generate(ctx, seq(1000)...), workers, 10*time.Millisecond)...)

	// Read a single result, then cancel and stop reading promptly.
	select {
	case <-merged:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for first result")
	}
	cancel()

	// After cancel, each worker can still hand over at most the result it
	// already finished, and each Merge forwarder at most the one it holds,
	// so no more than 2 results per worker may arrive.
	results := collect(t, merged)
	if max := 2 * workers; len(results) > max {
		t.Fatalf("expected cancellation to stop the pipeline promptly, got %d more results (max %d)", len(results), max)
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
