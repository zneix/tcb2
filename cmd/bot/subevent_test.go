package main

import "testing"

func TestPlaceholderReplacement(t *testing.T) {
	var testCases = []struct {
		format         string
		value          string
		login          string
		expectedOutput string
	}{
		{
			"NEW GAME 👉 {value} 👉 foo bar baz",
			"Artifact",
			"forsen",
			".me NEW GAME 👉 Artifact 👉 foo bar baz",
		},
		{
			"KKool GuitarTime {login} has gone live KKool GuitarTime",
			"Just Stalling",
			"forsen",
			".me KKool GuitarTime forsen has gone live KKool GuitarTime",
		},
		{
			"KKool GuitarTime {login} has gone live KKool GuitarTime",
			"Just Stalling",
			"{value}",
			".me KKool GuitarTime {value} has gone live KKool GuitarTime",
		},
		{
			"NEW GAME 👉 {value} foo",
			"this game has {login} in name",
			"zneix",
			".me NEW GAME 👉 this game has {login} in name foo",
		},
	}

	for _, test := range testCases {
		if output := createMessagePrefix(test.format, test.value, test.login); output != test.expectedOutput {
			t.Errorf("Expected %q, but resulted in %q", test.expectedOutput, output)
		}
	}
}
