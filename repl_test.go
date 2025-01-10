package main

import "testing"

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
			input:    "SOME thingy  bla  ",
			expected: []string{"some", "thingy", "bla"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "abc",
			expected: []string{"abc"},
		},
		{
			input:    "abc        ",
			expected: []string{"abc"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice
		if len(actual) != len(c.expected) {
			t.Errorf("Expected %d elements, got %d", len(c.expected), len(actual))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			if word != expectedWord {
				t.Errorf("Expected: %s, got %s", expectedWord, word)
			}
		}
	}
}
