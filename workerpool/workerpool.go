// Package workerpool runs tasks on a fixed number of goroutines fed by a
// bounded queue.
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
	// ErrPoolClosed is returned by Submit when the pool is closed or closing;
	// see Submit for the cases around Close.
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

// New starts workers goroutines with a queue of queueSize pending tasks.
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

// Submit queues t, blocking while the queue is full until there is room or
// ctx is done. It returns ErrNilTask for a nil task and ctx.Err() if ctx is
// already done. Otherwise, once Close has returned, it returns ErrPoolClosed.
// While Close is in progress, Submit may return ErrPoolClosed or ctx.Err(),
// or still queue its task (returning nil, and the task runs before Close
// returns), whether or not the queue is full.
func (p *Pool) Submit(ctx context.Context, t Task) error {
	if t == nil {
		return ErrNilTask
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	// The read lock stops Close from closing p.tasks mid-send.
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

// Close stops accepting tasks, waits for queued and running ones, and returns
// their errors joined. It must not be called from inside a task.
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

var errTaskGoexit = errors.New("workerpool: task called runtime.Goexit")

func (p *Pool) worker() {
	finished := false
	defer func() {
		if !finished {
			// A task called runtime.Goexit: replace this worker.
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
