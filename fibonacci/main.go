// Command fibonacci prints the first N terms of the Fibonacci series.
//
// Usage:
//
//	go run ./fibonacci        # prints the first 10 terms
//	go run ./fibonacci 15     # prints the first 15 terms
package main

import (
	"fmt"
	"os"
	"strconv"
)

// defaultCount is the number of terms printed when no argument is given.
const defaultCount = 10

// fib returns the first n terms of the Fibonacci series (starting at 0).
// It computes the series iteratively to avoid recursion overhead.
// For n <= 0 an empty slice is returned.
//
// Note: values are stored as uint64, so terms beyond fib(93) will overflow.
func fib(n int) []uint64 {
	if n <= 0 {
		return []uint64{}
	}

	series := make([]uint64, n)
	for i := 0; i < n; i++ {
		switch i {
		case 0:
			series[i] = 0
		case 1:
			series[i] = 1
		default:
			series[i] = series[i-1] + series[i-2]
		}
	}
	return series
}

func main() {
	count := defaultCount

	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err != nil || n < 0 {
			fmt.Fprintf(os.Stderr, "usage: %s [count]\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "count must be a non-negative integer\n")
			os.Exit(1)
		}
		count = n
	}

	series := fib(count)
	for i, v := range series {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(v)
	}
	fmt.Println()
}
