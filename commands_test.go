package main

import "testing"

func TestCatchRateDecreasesWithBaseExperience(t *testing.T) {
	tests := []struct {
		name       string
		baseExpLow int
		baseExpHi  int
	}{
		{name: "small increase", baseExpLow: 50, baseExpHi: 100},
		{name: "medium increase", baseExpLow: 100, baseExpHi: 200},
		{name: "large increase", baseExpLow: 200, baseExpHi: 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			low := catchRate(tt.baseExpLow)
			high := catchRate(tt.baseExpHi)
			if low <= high {
				t.Fatalf("catchRate(%d) = %f, catchRate(%d) = %f, want lower base experience to be easier to catch", tt.baseExpLow, low, tt.baseExpHi, high)
			}
		})
	}
}

func TestCatchRateBounds(t *testing.T) {
	if got := catchRate(0); got != 0.95 {
		t.Fatalf("catchRate(0) = %f, want 0.95", got)
	}

	if got := catchRate(10000); got != 0.05 {
		t.Fatalf("catchRate(10000) = %f, want 0.05", got)
	}
}
