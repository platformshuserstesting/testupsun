package main

import (
	"fmt"
	"strings"
)

// Fibonacci returns the first n terms of the Fibonacci series.
// The sequence starts with 0 and 1. For n <= 0 an empty slice is
// returned. The computation is iterative, running in O(n) time.
func Fibonacci(n int) []int {
	if n <= 0 {
		return []int{}
	}

	series := make([]int, n)
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
	const terms = 10

	series := Fibonacci(terms)

	parts := make([]string, len(series))
	for i, v := range series {
		parts[i] = fmt.Sprintf("%d", v)
	}

	fmt.Printf("First %d Fibonacci numbers: %s\n", terms, strings.Join(parts, " "))
}
