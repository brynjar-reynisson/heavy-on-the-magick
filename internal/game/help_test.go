package game

import (
	"strings"
	"testing"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
)

func TestHandleNameReportsPlayerName(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("NAME"))
	if !strings.Contains(got, "Axil") {
		t.Errorf("Handle(NAME) = %q, want it to mention Axil", got)
	}
}

// TestHandleHelpShowsRealHintScreen covers the real, disassembled,
// screenshot-cross-confirmed "SOME ADVICE" in-game hint screen (see
// help's doc comment) — the game's own actual text, not invented.
func TestHandleHelpShowsRealHintScreen(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("HELP"))
	for _, want := range []string{"SOME ADVICE", "APEX, THANKS", "GUARDS, DOOR", "BLAST without an object"} {
		if !strings.Contains(got, want) {
			t.Errorf("Handle(HELP) = %q, want it to contain %q", got, want)
		}
	}
}

func TestHandleGradeReportsCurrentGrade(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("GRADE"))
	if !strings.Contains(got, "Neophyte") {
		t.Errorf("Handle(GRADE) for a new player = %q, want it to mention Neophyte", got)
	}
}

// TestHandleSpellsListsRealSpells pins round 118: INVOKE joined this
// list because the manual's own explicit "Spells:" heading groups
// exactly I (Invoke), B (Blast), F (Freeze) - INVOKE was missing from
// this port's SPELLS listing entirely before, a real sourced omission,
// not just an incomplete aggregation.
func TestHandleSpellsListsRealSpells(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("SPELLS"))
	for _, want := range []string{"INVOKE", "BLAST", "FREEZE", "TRANSFUSION", "CALL"} {
		if !strings.Contains(got, want) {
			t.Errorf("Handle(SPELLS) = %q, want it to list %q", got, want)
		}
	}
}
