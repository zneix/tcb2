package main

import "testing"

func TestPlaceholderReplacement(t *testing.T) {
	var testCases = []struct {
		format         string
		value          string
		login          string
		redirect       bool
		expectedOutput string
	}{
		{
			"NEW GAME 👉 {value} 👉 foo bar baz",
			"Artifact",
			"forsen",
			false,
			".me NEW GAME 👉 Artifact 👉 foo bar baz",
		},
		{
			"KKool GuitarTime {login} has gone live KKool GuitarTime",
			"Just Stalling",
			"forsen",
			false,
			".me KKool GuitarTime forsen has gone live KKool GuitarTime",
		},
		{
			"KKool GuitarTime {login} has gone live KKool GuitarTime",
			"Just Stalling",
			"{value}",
			false,
			".me KKool GuitarTime {value} has gone live KKool GuitarTime",
		},
		{
			"NEW GAME 👉 {value} foo",
			"this game has {login} in name",
			"zneix",
			false,
			".me NEW GAME 👉 this game has {login} in name foo",
		},
		{
			"NEW GAME 👉 {value} 👉 whykingr ks romeo zulul",
			"PUBG: BATTLEGROUNDS",
			"forsen",
			true,
			".me [#forsen] NEW GAME 👉 PUBG: BATTLEGROUNDS 👉 whykingr ks romeo zulul",
		},
	}

	for _, test := range testCases {
		if output := createMessagePrefix(test.format, test.value, test.login, test.redirect); output != test.expectedOutput {
			t.Errorf("Expected %q, but resulted in %q", test.expectedOutput, output)
		}
	}
}
