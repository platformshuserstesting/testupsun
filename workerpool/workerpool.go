// Package workerpool provides a fixed-size pool of worker goroutines that
// execute tasks pulled from a bounded queue.
//
// Typical usage:
//
//	pool, err := workerpool.New(4, 16)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for i := 0; i < 100; i++ {
//		if err := pool.Submit(ctx, func() error {
//			// do work
//			return nil
//		}); err != nil {
//			log.Println(err)
//		}
//	}
//	if err := pool.Close(); err != nil {
//		log.Println("some tasks failed:", err)
//	}
//
// Submit blocks while the queue is full (backpressure) until there is room or
// the supplied context is done. Close stops accepting new tasks, waits for all
// queued tasks to finish and returns the errors reported by tasks, joined with
// errors.Join. Panics inside tasks are recovered and reported as errors.
package workerpool

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Task is a unit of work executed by a pool worker.
type Task func() error

var (
	// ErrInvalidWorkers is returned by New when the worker count is below 1.
	ErrInvalidWorkers = errors.New("workerpool: workers must be at least 1")
	// ErrInvalidQueueSize is returned by New when the queue size is negative.
	ErrInvalidQueueSize = errors.New("workerpool: queue size must not be negative")
	// ErrPoolClosed is returned by Submit after Close has been called.
	ErrPoolClosed = errors.New("workerpool: pool is closed")
	// ErrNilTask is returned by Submit when the task is nil.
	ErrNilTask = errors.New("workerpool: task is nil")
)

// Pool is a fixed-size pool of worker goroutines.
type Pool struct {
	workers int
	tasks   chan Task
	done    chan struct{} // closed when Close starts; unblocks pending Submits
	wg      sync.WaitGroup

	mu     sync.RWMutex // guards closed and the closing of tasks
	closed bool

	errMu sync.Mutex
	errs  []error

	closeOnce sync.Once
	closeErr  error
}

// New starts a pool with the given number of workers and a task queue that
// buffers up to queueSize pending tasks. A queueSize of 0 means Submit hands
// tasks directly to an idle worker.
func New(workers, queueSize int) (*Pool, error) {
	if workers < 1 {
		return nil, ErrInvalidWorkers
	}
	if queueSize < 0 {
		return nil, ErrInvalidQueueSize
	}
	p := &Pool{
		workers: workers,
		tasks:   make(chan Task, queueSize),
		done:    make(chan struct{}),
	}
	p.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go p.worker()
	}
	return p, nil
}

// Workers returns the number of worker goroutines in the pool.
func (p *Pool) Workers() int {
	return p.workers
}

// Submit queues t for execution. It blocks while the queue is full until
// there is room or ctx is done, in which case ctx.Err() is returned and the
// task is not run. If ctx is already done when Submit is called, the task is
// never queued. It returns ErrPoolClosed once Close has been called. A Submit
// that was waiting for queue space when Close started is released: it
// returns ErrPoolClosed, unless queue space opened up at the same moment, in
// which case the task is queued, runs before Close returns, and Submit
// returns nil.
func (p *Pool) Submit(ctx context.Context, t Task) error {
	if t == nil {
		return ErrNilTask
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	// Holding the read lock across the send guarantees Close cannot close the
	// channel while a send is in flight. Close first closes p.done, which
	// releases any Submit blocked on a full queue, so Close's write lock is
	// never held up by a waiting sender (even if tasks themselves call
	// Submit).
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.closed {
		return ErrPoolClosed
	}
	select {
	case p.tasks <- t:
		return nil
	case <-p.done:
		return ErrPoolClosed
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Close stops accepting new tasks, waits for all queued and running tasks to
// finish and returns the errors reported by tasks joined with errors.Join
// (nil if none failed). It is safe to call Close more than once; later calls
// return the same result.
//
// Close must not be called from inside a task: it waits for every worker to
// finish, including the one running the calling task, and would block forever.
func (p *Pool) Close() error {
	p.closeOnce.Do(func() {
		close(p.done)
		p.mu.Lock()
		p.closed = true
		close(p.tasks)
		p.mu.Unlock()

		p.wg.Wait()

		p.errMu.Lock()
		p.closeErr = errors.Join(p.errs...)
		p.errMu.Unlock()
	})
	return p.closeErr
}

// errTaskGoexit is reported when a task ends its goroutine with
// runtime.Goexit (for example t.FailNow called inside a task).
var errTaskGoexit = errors.New("workerpool: task called runtime.Goexit")

func (p *Pool) worker() {
	finished := false
	defer func() {
		if !finished {
			// Panics are recovered in run, so the only way to get here early
			// is a task calling runtime.Goexit. Report it and start a
			// replacement so the pool keeps its worker count. Add before Done
			// keeps the WaitGroup above zero, so Close keeps waiting.
			p.addErr(errTaskGoexit)
			p.wg.Add(1)
			go p.worker()
		}
		p.wg.Done()
	}()
	for t := range p.tasks {
		if err := p.run(t); err != nil {
			p.addErr(err)
		}
	}
	finished = true
}

func (p *Pool) addErr(err error) {
	p.errMu.Lock()
	p.errs = append(p.errs, err)
	p.errMu.Unlock()
}

func (p *Pool) run(t Task) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("workerpool: task panicked: %v", r)
		}
	}()
	return t()
}
