package main

import "testing"

func TestF(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 3},
		{3, 6},
		{4, 10},
		{5, 15},
		{10, 45},
	}

	for _, test := range tests {
		if got := f(test.input); got != test.want {
			t.Errorf("f(%d) = %d; want %d", test.input, got, test.want)
		}
	}
}
