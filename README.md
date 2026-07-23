# Fibonacci Sample (Go)

A small, self-contained Go program that computes and prints the Fibonacci
series. The series starts at 0 (`0, 1, 1, 2, 3, 5, 8, ...`) and uses
`math/big` so large counts do not overflow.

## Usage

```sh
# Print the first 10 Fibonacci numbers (default)
go run ./fibonacci

# Print the first 15 Fibonacci numbers
go run ./fibonacci 15
```

Passing a non-numeric or negative count prints a usage message and exits with a
non-zero status. A count of `0` prints an empty series.

## Testing

```sh
go test ./...
go vet ./...
```
