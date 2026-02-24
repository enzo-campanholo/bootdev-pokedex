package main

import (
	"slices"
	"testing"
)

type testCase struct {
	name     string
	input    string
	expected []string
}

func TestCleanInput(t *testing.T) {
	cases := []testCase{
		{
			name:     "empty",
			input:    "  ",
			expected: []string{},
		},
		{
			name:     "single word",
			input:    "  hello  ",
			expected: []string{"hello"},
		},
		{
			name:     "two words",
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			name:     "mixed case",
			input:    "  HellO  World  ",
			expected: []string{"hello", "world"},
		},
		{
			name:     "multiple words",
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := cleanInput(c.input)

			if !slices.Equal(actual, c.expected) {
				t.Fatalf("cleanInput(%q) = %v, want %v", c.input, actual, c.expected)
			}
		})
	}
}
