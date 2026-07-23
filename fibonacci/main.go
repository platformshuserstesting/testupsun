// Command fibonacci is a small sample program that computes and prints the
// Fibonacci series.
//
// Usage:
//
//	go run ./fibonacci        # prints the first 10 Fibonacci numbers
//	go run ./fibonacci 15     # prints the first 15 Fibonacci numbers
//
// The series starts at 0 (0, 1, 1, 2, 3, 5, 8, ...). math/big is used so that
// arbitrarily large counts do not overflow.
package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
)

// defaultCount is the number of terms printed when no argument is given.
const defaultCount = 10

// fib returns the first n terms of the Fibonacci series. For n <= 0 it returns
// an empty slice. It uses big.Int to support arbitrarily large values without
// overflow.
func fib(n int) []*big.Int {
	if n <= 0 {
		return []*big.Int{}
	}

	series := make([]*big.Int, n)
	a, b := big.NewInt(0), big.NewInt(1)
	for i := 0; i < n; i++ {
		series[i] = new(big.Int).Set(a)
		a, b = b, new(big.Int).Add(a, b)
	}
	return series
}

func main() {
	count := defaultCount

	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err != nil || n < 0 {
			fmt.Fprintf(os.Stderr, "error: count must be a non-negative integer\n")
			fmt.Fprintf(os.Stderr, "usage: %s [count]\n", os.Args[0])
			os.Exit(1)
		}
		count = n
	}

	for i, v := range fib(count) {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(v.String())
	}
	fmt.Println()
}
