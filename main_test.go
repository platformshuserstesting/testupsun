package main

import (
	"reflect"
	"testing"
)

func TestFibonacci(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want []int
	}{
		{name: "negative", n: -1, want: []int{}},
		{name: "zero", n: 0, want: []int{}},
		{name: "one", n: 1, want: []int{0}},
		{name: "two", n: 2, want: []int{0, 1}},
		{name: "five", n: 5, want: []int{0, 1, 1, 2, 3}},
		{name: "ten", n: 10, want: []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Fibonacci(tt.n)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fibonacci(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}
