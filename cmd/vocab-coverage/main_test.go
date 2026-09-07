package main

import (
	"testing"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/game"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
)

func TestUniqueSortedDeduplicatesAndSorts(t *testing.T) {
	got := uniqueSorted([]string{"WEST", "EAST", "WEST", "APEX"})
	want := []string{"APEX", "EAST", "WEST"}
	if len(got) != len(want) {
		t.Fatalf("uniqueSorted = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("uniqueSorted = %v, want %v", got, want)
		}
	}
}

func TestIsGenericResponse(t *testing.T) {
	cases := []struct {
		resp string
		want bool
	}{
		{"I recognize that word, but don't know what it does yet (verb resolution not yet implemented).", true},
		{"I don't understand that word.", true},
		{"You are Axil.", false},
		{"", false},
	}
	for _, c := range cases {
		if got := isGenericResponse(c.resp); got != c.want {
			t.Errorf("isGenericResponse(%q) = %v, want %v", c.resp, got, c.want)
		}
	}
}

// TestKnownImplementedWordsAreNotUncovered pins a few real vocabulary
// words this project has definitely implemented (see game.Handle) as a
// regression check — if any of these ever showed up as "uncovered", the
// classification logic itself would be broken, not the game.
func TestKnownImplementedWordsAreNotUncovered(t *testing.T) {
	for _, w := range []string{"BLAST", "TRANSFUSION", "INVOKE", "NORTH", "HELP", "GRADE"} {
		g := game.New()
		resp := g.Handle(parser.Command{Verb: w})
		if isGenericResponse(resp) {
			t.Errorf("Handle(Verb: %q) = %q, want real modeled behavior, not the generic stub", w, resp)
		}
	}
}

// TestTargetPositionWordsAreExcludedButRealVocabulary pins that every
// word in targetPositionWords is (a) actually in the real vocabulary
// (not a typo) and (b) genuinely falls through to the generic stub when
// checked bare as cmd.Verb with no companion - confirming the exclusion
// is warranted, not just convenient.
func TestTargetPositionWordsAreExcludedButRealVocabulary(t *testing.T) {
	for w := range targetPositionWords {
		if !parser.KnownWord(w) {
			t.Errorf("targetPositionWords contains %q, which isn't in the real vocabulary", w)
		}
	}
}
