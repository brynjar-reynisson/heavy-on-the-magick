package game

import (
	"strings"
	"testing"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
)

func TestHandleSwapRecognized(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("Z")) // "Z" expands to SWAP
	if strings.Contains(got, "don't understand") {
		t.Errorf("Handle(Z) = %q, want a recognized (if stubbed) response, not the unknown-word rejection", got)
	}
}

// TestHandleLeftRightRecognized covers LEFT/RIGHT (Merphish keywords
// L/R) - real, frequently-used CASA walkthrough commands, recognized
// as real (not lumped in with "I don't understand that word") but
// honestly stubbed since this port has no facing-direction state.
func TestHandleLeftRightRecognized(t *testing.T) {
	for _, keyword := range []string{"L", "R"} {
		g := New()
		got := g.Handle(parser.Parse(keyword))
		if strings.Contains(got, "don't understand") {
			t.Errorf("Handle(%s) = %q, want it recognized as a real command, not an unknown word", keyword, got)
		}
	}
}

func TestHandleOptionsShowsMenu(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("OPTIONS"))
	if !strings.Contains(got, "Realign Status") {
		t.Errorf("Handle(OPTIONS) = %q, want it to list the confirmed menu including Realign Status", got)
	}
}

func TestHandleOptionsRealignRerollsStats(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("O REALIGN STATUS"))
	if !strings.Contains(got, "Realign Status") {
		t.Errorf("Handle(O REALIGN STATUS) = %q, want it to confirm the realign", got)
	}
	if g.Player.Stamina <= 0 || g.Player.Skill <= 0 || g.Player.Luck <= 0 {
		t.Errorf("Player after Realign = %+v, want all stats positive", g.Player)
	}
}

// TestHandleHalt covers HALT (Merphish "H") - round 85 corrected the
// response to reflect the manual's precise definition ("abandon the
// command being actioned and the rest of any outstanding command
// string"), not just a generic "Halted."
func TestHandleHalt(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("HALT"))
	if !strings.Contains(got, "halted") {
		t.Errorf("Handle(HALT) = %q, want it to acknowledge the command was halted", got)
	}
}
