package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "  ",
			expected: []string{},
		},
		{
			input:    "pokemon",
			expected: []string{"pokemon"},
		},
		{
			input:    "  HELLO    WORLD  ",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		// Check the length of the actual slice against the expected slice
		if len(actual) != len(c.expected) {
			t.Errorf("For input '%s', expected length %d, but got %d",
				c.input, len(c.expected), len(actual))
			continue
		}

		// Check each word in the slice
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("For input '%s' at index %d, expected '%s', but got '%s'",
					c.input, i, expectedWord, word)
			}
		}
	}
}
