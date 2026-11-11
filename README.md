# testupsun

A small Go worker pool: a fixed number of goroutines ("threads") pull tasks
from a bounded queue.

- `workerpool/` — the reusable `workerpool` package.
- `cmd/workerpool/` — a demo program that squares numbers on the pool.

## Usage

```go
pool, err := workerpool.New(4, 16) // 4 workers, queue of 16 pending tasks
if err != nil {
	log.Fatal(err)
}
err = pool.Submit(ctx, func() error {
	// do work
	return nil
})
// ...
if err := pool.Close(); err != nil { // waits for queued tasks to finish
	log.Println("some tasks failed:", err)
}
```

- `Submit` blocks while the queue is full until there is room or `ctx` is done.
- `Submit` returns `ErrPoolClosed` after `Close`, and `ErrNilTask` for a nil task.
- `Close` waits for all queued tasks to finish and returns the errors from
  tasks (joined with `errors.Join`). Panics in tasks are recovered and
  reported as errors, as is a task that calls `runtime.Goexit` (the worker is
  replaced). Calling `Close` more than once is safe. A `Submit` still waiting
  for queue space when `Close` starts is released with `ErrPoolClosed`,
  unless space opened up at the same moment; then the task is queued, runs,
  and `Submit` returns nil.
  Do not call `Close` from inside a task, because it would wait on itself.

## Run the demo

```sh
go run ./cmd/workerpool -workers 4 -jobs 20
```

Flags: `-workers` (default 4), `-jobs` (default 20), `-queue` (default 10).

## Run the tests

```sh
go test -race ./...
```
