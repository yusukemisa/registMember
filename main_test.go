package main

import "testing"

func Test_Sample(t *testing.T) {
	tests := []struct {
		caseName string
		in int
		out int
	}{
		{
			caseName:"sample test 1",
			in: 1,
			out: 1,
		},
		{
			caseName:"sample test 2",
			in: 2,
			out: 2,
		},
	}

	for _,tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			if tt.in != tt.out {
				t.Errorf("want=%v,got=%v",tt.in,tt.out)
			}
		})
	}
}
