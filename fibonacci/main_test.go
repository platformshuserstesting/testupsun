package main

import (
	"math/big"
	"testing"
)

func TestFib(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want []int64
	}{
		{name: "negative", n: -3, want: []int64{}},
		{name: "zero", n: 0, want: []int64{}},
		{name: "one", n: 1, want: []int64{0}},
		{name: "two", n: 2, want: []int64{0, 1}},
		{name: "ten", n: 10, want: []int64{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fib(tt.n)
			if len(got) != len(tt.want) {
				t.Fatalf("fib(%d) returned %d terms, want %d", tt.n, len(got), len(tt.want))
			}
			for i, v := range got {
				want := big.NewInt(tt.want[i])
				if v.Cmp(want) != 0 {
					t.Errorf("fib(%d)[%d] = %s, want %s", tt.n, i, v.String(), want.String())
				}
			}
		})
	}
}

func TestFibLargeNoOverflow(t *testing.T) {
	// The 100th term (index 99) is 218922995834555169026, which overflows
	// uint64/int64; math/big must handle it correctly.
	got := fib(100)
	want, ok := new(big.Int).SetString("218922995834555169026", 10)
	if !ok {
		t.Fatal("failed to parse expected big.Int")
	}
	if got[99].Cmp(want) != 0 {
		t.Errorf("fib(100)[99] = %s, want %s", got[99].String(), want.String())
	}
}
