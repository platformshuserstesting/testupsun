# fibonacci-sample

A small, self-contained Go sample that computes and prints the Fibonacci series.

## Running

```sh
# Print the first 10 terms (default)
go run ./fibonacci

# Print the first 15 terms
go run ./fibonacci 15
```

The count argument must be a non-negative integer. Passing `0` prints an
empty series. Values are stored as `uint64`, so terms beyond `fib(93)`
will overflow.

## Testing

```sh
go vet ./...
go test ./...
```
