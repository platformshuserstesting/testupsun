# Go fan-out / fan-in demo

A small Go program that demonstrates the **fan-out / fan-in** pattern with
goroutines and channels, using only the standard library.

```
              +-> worker 1 --+
generate --> -+-> worker 2 --+-> merge --> main (prints results)
              +-> worker N --+
```

- **Fan-out:** `fanio.Generate` sends jobs on one channel. `fanio.FanOut`
  starts N workers that all read from that same channel, so each job is
  picked up by exactly one worker and the work runs concurrently.
- **Fan-in:** each worker writes to its own output channel. `fanio.Merge`
  forwards every worker channel into a single channel and closes it (using a
  `sync.WaitGroup`) only once all workers are done.
- Every stage selects on `ctx.Done()`, so Ctrl-C stops the pipeline cleanly
  without leaking goroutines.

The "work" is squaring a number after a short simulated delay.

## Run

Requires Go 1.22 or newer.

```sh
go run ./cmd/fanin-fanout
go run ./cmd/fanin-fanout -workers 3 -jobs 10 -delay 100ms
```

Flags:

| Flag       | Default | Description                         |
|------------|---------|-------------------------------------|
| `-workers` | `4`     | number of worker goroutines (1-100) |
| `-jobs`    | `20`    | number of jobs to generate          |
| `-delay`   | `50ms`  | simulated work time per job         |

Results are printed as they arrive, so their order changes between runs. A
summary at the end shows the total, the sum of outputs, and how many jobs each
worker handled. If interrupted with Ctrl-C, the program prints the partial
summary and exits with status 130.

## Test

```sh
go vet ./...
go test -race ./...
```
