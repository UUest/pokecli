package main

import (
	"testing";
	"github.com/UUest/pokecli/cleaninput"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input: "heya erf",
			expected: []string{"heya", "erf"},
		},
	}

	for _, c := range cases {
		actual := cleaninput.CleanInput(c.input)
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Expected %s, got %s", expectedWord, word)
			}
		}
	}
}
