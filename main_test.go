package main

import "testing"

// -----------------------------------------------------------------------
// levenshtein
// -----------------------------------------------------------------------

func TestLevenshtein_Identical(t *testing.T) {
	if d := levenshtein("daemon", "daemon"); d != 0 {
		t.Errorf("expected 0, got %d", d)
	}
}

func TestLevenshtein_EmptyStrings(t *testing.T) {
	if d := levenshtein("", ""); d != 0 {
		t.Errorf("expected 0, got %d", d)
	}
}

func TestLevenshtein_OneEmpty(t *testing.T) {
	if d := levenshtein("abc", ""); d != 3 {
		t.Errorf("expected 3, got %d", d)
	}
	if d := levenshtein("", "abc"); d != 3 {
		t.Errorf("expected 3, got %d", d)
	}
}

func TestLevenshtein_SingleSubstitution(t *testing.T) {
	if d := levenshtein("list", "lust"); d != 1 {
		t.Errorf("expected 1, got %d", d)
	}
}

func TestLevenshtein_SingleInsertion(t *testing.T) {
	if d := levenshtein("copy", "coppy"); d != 1 {
		t.Errorf("expected 1, got %d", d)
	}
}

func TestLevenshtein_SingleDeletion(t *testing.T) {
	if d := levenshtein("clear", "cler"); d != 1 {
		t.Errorf("expected 1, got %d", d)
	}
}

func TestLevenshtein_CompletelyDifferent(t *testing.T) {
	d := levenshtein("abc", "xyz")
	if d != 3 {
		t.Errorf("expected 3, got %d", d)
	}
}

// -----------------------------------------------------------------------
// suggest
// -----------------------------------------------------------------------

func TestSuggest_ExactMatch(t *testing.T) {
	if got := suggest("daemon"); got != "daemon" {
		t.Errorf("expected %q, got %q", "daemon", got)
	}
}

func TestSuggest_OneTypo(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"daemom", "daemon"},
		{"stat us", "status"},
		{"lst", "list"},
		{"cler", "clear"},
		{"upgade", "upgrade"},
	}
	for _, tc := range cases {
		if got := suggest(tc.input); got != tc.want {
			t.Errorf("suggest(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestSuggest_TooFarOff(t *testing.T) {
	// Input with edit distance > 3 should return no suggestion.
	if got := suggest("xxxxxxxxxx"); got != "" {
		t.Errorf("expected empty suggestion for gibberish, got %q", got)
	}
}

func TestSuggest_EmptyInput(t *testing.T) {
	// All known commands have length > 3, so distance from "" is their length.
	// Shortest known command is "run" (3) and "pin" (3) — distance 3, which is
	// exactly the threshold (< 4), so one of them may be returned.
	// Just assert it doesn't panic and returns a string.
	_ = suggest("")
}
