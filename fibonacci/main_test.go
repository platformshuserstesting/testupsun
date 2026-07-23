package main

import (
	"reflect"
	"testing"
)

func TestFib(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want []uint64
	}{
		{name: "negative", n: -1, want: []uint64{}},
		{name: "zero", n: 0, want: []uint64{}},
		{name: "one", n: 1, want: []uint64{0}},
		{name: "two", n: 2, want: []uint64{0, 1}},
		{name: "ten", n: 10, want: []uint64{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fib(tt.n)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fib(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}
