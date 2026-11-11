package workerpool

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// guard is the maximum time any blocking step may take before the test fails.
// It is a failure guard only; tests never rely on it for correctness.
const guard = 5 * time.Second

// waitFor fails the test if ch does not become readable within guard.
func waitFor(t *testing.T, ch <-chan struct{}, what string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(guard):
		t.Fatalf("timed out waiting for %s", what)
	}
}

// closeWithin runs p.Close in a goroutine and returns its result, failing the
// test if Close does not return within guard.
func closeWithin(t *testing.T, p *Pool) error {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- p.Close() }()
	select {
	case err := <-done:
		return err
	case <-time.After(guard):
		t.Fatal("timed out waiting for Close")
		return nil
	}
}

func TestNewValidation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		workers   int
		queueSize int
		wantErr   error
	}{
		{"zero workers", 0, 1, ErrInvalidWorkers},
		{"negative workers", -3, 1, ErrInvalidWorkers},
		{"negative queue", 1, -1, ErrInvalidQueueSize},
		{"unbuffered queue", 1, 0, nil},
		{"valid", 4, 10, nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p, err := New(tc.workers, tc.queueSize)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("New(%d, %d) error = %v, want %v", tc.workers, tc.queueSize, err, tc.wantErr)
			}
			if tc.wantErr != nil {
				if p != nil {
					t.Fatalf("New returned non-nil pool with error")
				}
				return
			}
			if got := p.Workers(); got != tc.workers {
				t.Errorf("Workers() = %d, want %d", got, tc.workers)
			}
			if err := closeWithin(t, p); err != nil {
				t.Errorf("Close() = %v, want nil", err)
			}
		})
	}
}

func TestAllTasksRunExactlyOnce(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		workers   int
		queueSize int
		tasks     int
	}{
		{"single worker unbuffered", 1, 0, 50},
		{"single worker buffered", 1, 10, 50},
		{"many workers unbuffered", 8, 0, 500},
		{"many workers buffered", 8, 32, 500},
		{"more workers than tasks", 16, 4, 5},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p, err := New(tc.workers, tc.queueSize)
			if err != nil {
				t.Fatal(err)
			}
			counts := make([]int32, tc.tasks)
			for i := 0; i < tc.tasks; i++ {
				if err := p.Submit(context.Background(), func() error {
					atomic.AddInt32(&counts[i], 1)
					return nil
				}); err != nil {
					t.Fatalf("Submit(%d) = %v", i, err)
				}
			}
			if err := closeWithin(t, p); err != nil {
				t.Fatalf("Close() = %v, want nil", err)
			}
			for i := range counts {
				if c := atomic.LoadInt32(&counts[i]); c != 1 {
					t.Errorf("task %d ran %d times, want 1", i, c)
				}
			}
		})
	}
}

func TestBoundedConcurrency(t *testing.T) {
	t.Parallel()
	const workers, tasks = 4, 200
	p, err := New(workers, 8)
	if err != nil {
		t.Fatal(err)
	}
	var inFlight, maxInFlight int32
	for i := 0; i < tasks; i++ {
		if err := p.Submit(context.Background(), func() error {
			n := atomic.AddInt32(&inFlight, 1)
			for {
				m := atomic.LoadInt32(&maxInFlight)
				if n <= m || atomic.CompareAndSwapInt32(&maxInFlight, m, n) {
					break
				}
			}
			atomic.AddInt32(&inFlight, -1)
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	if err := closeWithin(t, p); err != nil {
		t.Fatal(err)
	}
	if m := atomic.LoadInt32(&maxInFlight); m > workers || m < 1 {
		t.Fatalf("max in-flight tasks = %d, want between 1 and %d", m, workers)
	}
}

func TestWorkersRunInParallel(t *testing.T) {
	t.Parallel()
	const workers = 5
	p, err := New(workers, 0)
	if err != nil {
		t.Fatal(err)
	}
	var started sync.WaitGroup
	started.Add(workers)
	release := make(chan struct{})
	for i := 0; i < workers; i++ {
		if err := p.Submit(context.Background(), func() error {
			started.Done()
			<-release
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	allStarted := make(chan struct{})
	go func() {
		started.Wait()
		close(allStarted)
	}()
	waitFor(t, allStarted, "all workers to start concurrently")
	close(release)
	if err := closeWithin(t, p); err != nil {
		t.Fatal(err)
	}
}

func TestSubmitNilTask(t *testing.T) {
	t.Parallel()
	p, err := New(1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Submit(context.Background(), nil); !errors.Is(err, ErrNilTask) {
		t.Fatalf("Submit(nil) = %v, want %v", err, ErrNilTask)
	}
	if err := closeWithin(t, p); err != nil {
		t.Fatal(err)
	}
}

func TestSubmitAfterClose(t *testing.T) {
	t.Parallel()
	p, err := New(2, 2)
	if err != nil {
		t.Fatal(err)
	}
	if err := closeWithin(t, p); err != nil {
		t.Fatal(err)
	}
	var ran atomic.Bool
	err = p.Submit(context.Background(), func() error {
		ran.Store(true)
		return nil
	})
	if !errors.Is(err, ErrPoolClosed) {
		t.Fatalf("Submit after Close = %v, want %v", err, ErrPoolClosed)
	}
	if ran.Load() {
		t.Fatal("task submitted after Close ran")
	}
}

func TestSubmitBackpressureContext(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		makeCtx func() (context.Context, context.CancelFunc)
		wantErr error
	}{
		{
			name: "deadline exceeded",
			makeCtx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 20*time.Millisecond)
			},
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "canceled",
			makeCtx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx, cancel
			},
			wantErr: context.Canceled,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p, err := New(1, 1)
			if err != nil {
				t.Fatal(err)
			}
			release := make(chan struct{})
			started := make(chan struct{})
			// Occupy the only worker.
			if err := p.Submit(context.Background(), func() error {
				close(started)
				<-release
				return nil
			}); err != nil {
				t.Fatal(err)
			}
			waitFor(t, started, "blocking task to start")
			// Fill the queue.
			if err := p.Submit(context.Background(), func() error { return nil }); err != nil {
				t.Fatal(err)
			}

			ctx, cancel := tc.makeCtx()
			defer cancel()
			var ran atomic.Bool
			err = p.Submit(ctx, func() error {
				ran.Store(true)
				return nil
			})
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Submit on full queue = %v, want %v", err, tc.wantErr)
			}
			close(release)
			if err := closeWithin(t, p); err != nil {
				t.Fatal(err)
			}
			if ran.Load() {
				t.Fatal("task rejected by context was run")
			}
		})
	}
}

func TestCloseDrainsQueuedTasks(t *testing.T) {
	t.Parallel()
	const queued = 10
	p, err := New(1, queued)
	if err != nil {
		t.Fatal(err)
	}
	release := make(chan struct{})
	started := make(chan struct{})
	if err := p.Submit(context.Background(), func() error {
		close(started)
		<-release
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	waitFor(t, started, "blocking task to start")
	var done int32
	for i := 0; i < queued; i++ {
		if err := p.Submit(context.Background(), func() error {
			atomic.AddInt32(&done, 1)
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	closed := make(chan error, 1)
	go func() { closed <- p.Close() }()

	// Close must wait for the blocked task and the queue to drain.
	select {
	case err := <-closed:
		t.Fatalf("Close returned early with %v while a task was still running", err)
	default:
	}
	close(release)
	select {
	case err := <-closed:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(guard):
		t.Fatal("timed out waiting for Close")
	}
	if got := atomic.LoadInt32(&done); got != queued {
		t.Fatalf("%d queued tasks ran before Close returned, want %d", got, queued)
	}
}

func TestTaskErrorsAreAggregated(t *testing.T) {
	t.Parallel()
	errA := errors.New("task A failed")
	errB := errors.New("task B failed")
	p, err := New(3, 10)
	if err != nil {
		t.Fatal(err)
	}
	var ok int32
	for _, e := range []error{errA, nil, errB, nil, errA} {
		if err := p.Submit(context.Background(), func() error {
			if e == nil {
				atomic.AddInt32(&ok, 1)
			}
			return e
		}); err != nil {
			t.Fatal(err)
		}
	}
	err = closeWithin(t, p)
	if err == nil {
		t.Fatal("Close() = nil, want aggregated error")
	}
	if !errors.Is(err, errA) || !errors.Is(err, errB) {
		t.Fatalf("Close() = %v, want it to wrap %v and %v", err, errA, errB)
	}
	if n := len(err.(interface{ Unwrap() []error }).Unwrap()); n != 3 {
		t.Fatalf("Close() wrapped %d errors, want 3", n)
	}
	if got := atomic.LoadInt32(&ok); got != 2 {
		t.Fatalf("%d successful tasks ran, want 2", got)
	}
}

func TestTaskPanicIsRecovered(t *testing.T) {
	t.Parallel()
	p, err := New(1, 10)
	if err != nil {
		t.Fatal(err)
	}
	var after int32
	if err := p.Submit(context.Background(), func() error { panic("boom") }); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		if err := p.Submit(context.Background(), func() error {
			atomic.AddInt32(&after, 1)
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	err = closeWithin(t, p)
	if err == nil || !strings.Contains(err.Error(), "task panicked: boom") {
		t.Fatalf("Close() = %v, want panic error", err)
	}
	if got := atomic.LoadInt32(&after); got != 5 {
		t.Fatalf("%d tasks ran after the panic, want 5", got)
	}
}

func TestTaskGoexitReplacesWorker(t *testing.T) {
	t.Parallel()
	p, err := New(1, 10)
	if err != nil {
		t.Fatal(err)
	}
	var after int32
	if err := p.Submit(context.Background(), func() error {
		runtime.Goexit()
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		if err := p.Submit(context.Background(), func() error {
			atomic.AddInt32(&after, 1)
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	err = closeWithin(t, p)
	if !errors.Is(err, errTaskGoexit) {
		t.Fatalf("Close() = %v, want %v", err, errTaskGoexit)
	}
	if got := atomic.LoadInt32(&after); got != 5 {
		t.Fatalf("%d tasks ran after the Goexit, want 5", got)
	}
}

func TestSubmitWithDoneContextNeverQueues(t *testing.T) {
	t.Parallel()
	p, err := New(1, 10) // plenty of queue space
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var ran atomic.Bool
	for i := 0; i < 100; i++ {
		if err := p.Submit(ctx, func() error {
			ran.Store(true)
			return nil
		}); !errors.Is(err, context.Canceled) {
			t.Fatalf("Submit with canceled ctx = %v, want %v", err, context.Canceled)
		}
	}
	if err := closeWithin(t, p); err != nil {
		t.Fatal(err)
	}
	if ran.Load() {
		t.Fatal("task submitted with a canceled context was run")
	}
}

func TestCloseIsIdempotent(t *testing.T) {
	t.Parallel()
	sentinel := errors.New("failed")
	p, err := New(2, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Submit(context.Background(), func() error { return sentinel }); err != nil {
		t.Fatal(err)
	}
	first := closeWithin(t, p)
	second := closeWithin(t, p)
	if !errors.Is(first, sentinel) || !errors.Is(second, sentinel) {
		t.Fatalf("Close() = %v then %v, want both to wrap %v", first, second, sentinel)
	}
}

func TestConcurrentSubmitAndClose(t *testing.T) {
	t.Parallel()
	p, err := New(4, 4)
	if err != nil {
		t.Fatal(err)
	}
	const submitters, perSubmitter = 16, 50
	var accepted, ran int32
	var wg sync.WaitGroup
	start := make(chan struct{})
	errCh := make(chan error, submitters*perSubmitter)
	for i := 0; i < submitters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < perSubmitter; j++ {
				err := p.Submit(context.Background(), func() error {
					atomic.AddInt32(&ran, 1)
					return nil
				})
				switch {
				case err == nil:
					atomic.AddInt32(&accepted, 1)
				case errors.Is(err, ErrPoolClosed):
				default:
					errCh <- err
				}
			}
		}()
	}
	close(start)
	var closeErrs [2]error
	var cwg sync.WaitGroup
	for i := range closeErrs {
		cwg.Add(1)
		go func() {
			defer cwg.Done()
			closeErrs[i] = p.Close()
		}()
	}
	allDone := make(chan struct{})
	go func() {
		wg.Wait()
		cwg.Wait()
		close(allDone)
	}()
	waitFor(t, allDone, "submitters and Close to finish")
	close(errCh)
	for err := range errCh {
		t.Errorf("unexpected Submit error: %v", err)
	}
	for i, err := range closeErrs {
		if err != nil {
			t.Errorf("Close #%d = %v, want nil", i, err)
		}
	}
	if a, r := atomic.LoadInt32(&accepted), atomic.LoadInt32(&ran); a != r {
		t.Fatalf("accepted %d tasks but %d ran", a, r)
	}
}

// A Submit blocked on a full queue plus a task that itself calls Submit must
// not deadlock Close.
func TestCloseWithBlockedAndNestedSubmit(t *testing.T) {
	t.Parallel()
	p, err := New(1, 1)
	if err != nil {
		t.Fatal(err)
	}
	started := make(chan struct{})
	gate := make(chan struct{})
	nestedErr := make(chan error, 1)
	err = p.Submit(context.Background(), func() error {
		close(started)
		<-gate
		nestedErr <- p.Submit(context.Background(), func() error { return nil })
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	waitFor(t, started, "first task to start")
	// Fill the queue, then block an external Submit on it.
	if err := p.Submit(context.Background(), func() error { return nil }); err != nil {
		t.Fatal(err)
	}
	blockedErr := make(chan error, 1)
	go func() {
		blockedErr <- p.Submit(context.Background(), func() error { return nil })
	}()

	closeErr := make(chan error, 1)
	go func() { closeErr <- p.Close() }()
	close(gate)

	select {
	case err := <-closeErr:
		if err != nil {
			t.Fatalf("Close() = %v, want nil", err)
		}
	case <-time.After(guard):
		t.Fatal("timed out waiting for Close")
	}
	for name, ch := range map[string]chan error{"blocked": blockedErr, "nested": nestedErr} {
		select {
		case err := <-ch:
			if err != nil && !errors.Is(err, ErrPoolClosed) {
				t.Errorf("%s Submit = %v, want nil or %v", name, err, ErrPoolClosed)
			}
		case <-time.After(guard):
			t.Fatalf("timed out waiting for %s Submit", name)
		}
	}
}
