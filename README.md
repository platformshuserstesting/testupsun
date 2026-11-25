# Fan-out / fan-in with goroutines

A small Go program that demonstrates the fan-out / fan-in concurrency pattern.

```
                       +-> worker 1 --+
                       |              |
Generate --> jobs chan +-> worker 2 --+--> Merge --> results chan --> main
 (producer)            |              |   (fan-in)
                       +-> worker 3 --+
                          (fan-out)
```

- **Fan-out:** several worker goroutines read from the *same* jobs channel, so
  the work is shared between them.
- **Fan-in:** `Merge` starts one goroutine per worker result channel and
  forwards everything onto a single channel. A `sync.WaitGroup` closes that
  channel once every worker is done.

Every stage selects on `ctx.Done()`, so cancelling the context (or hitting a
timeout) shuts down the whole pipeline without leaking goroutines.

The pattern lives in the [`fanio`](fanio/fanio.go) package; `main.go` wires it
together.

## Run

```sh
go run .                                   # 3 workers, numbers 1..10, sum = 385
go run . -workers 5 -n 20                  # custom worker count and input size
go run . -n 20 -delay 10ms -timeout 30ms   # stop early on timeout
```

Flags:

| Flag       | Default | Description                                   |
|------------|---------|-----------------------------------------------|
| `-workers` | `3`     | number of worker goroutines (at most 10000)   |
| `-n`       | `10`    | process the numbers 1..n (at most 1000000)    |
| `-delay`   | `0`     | simulated work per item, e.g. `10ms`          |
| `-timeout` | `0`     | stop the pipeline after this long (0 = never) |

Results are printed in whatever order the workers finish, e.g.
`worker 2: 4 -> 16`.

## Test

```sh
go test -race ./...
```
