# Fibonacci Series (Go sample)

A small, self-contained Go program that computes and prints the Fibonacci
series. It exposes a reusable `Fibonacci(n int) []int` function backed by
table-driven unit tests.

## Requirements

- Go 1.22 or newer

## Run

```sh
go run .
```

Expected output:

```
First 10 Fibonacci numbers: 0 1 1 2 3 5 8 13 21 34
```

## Test

```sh
go test ./...
```
