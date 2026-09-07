package game

import (
	"strings"
	"testing"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
)

func TestHandleTalkToApex(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("APEX, TALK"))
	if !strings.Contains(got, "Apex") {
		t.Errorf("Handle(APEX, TALK) = %q, want it to recognize the real NPC Apex the Ogre", got)
	}
}

// TestHandleApexDoorWithNoHintFallsBackToTalk covers apexDoorHint's
// honest fallback: a room with no real, sourced DoorHints (like the
// starting room) gets the same generic response as bare "APEX, TALK",
// not a fabricated riddle.

// TestHandleApexDoorWithNoHintFallsBackToTalk covers apexDoorHint's
// honest fallback: a room with no real, sourced DoorHints (like the
// starting room) gets the same generic response as bare "APEX, TALK",
// not a fabricated riddle.
func TestHandleApexDoorWithNoHintFallsBackToTalk(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("APEX, DOOR"))
	if !strings.Contains(got, "Apex") {
		t.Errorf("Handle(APEX, DOOR) with no real hint here = %q, want the generic Apex response", got)
	}
	if strings.Contains(got, "riddle") {
		t.Errorf("Handle(APEX, DOOR) with no real hint here = %q, must not fabricate a riddle", got)
	}
}

// TestHandleApexDoorGivesRealWolfdorpHints covers round 159's real,
// sourced riddle content (see world.Room.DoorHints's doc comment) at
// Wolfdorp, an already-real, reachable CollodonsPile room.

// TestHandleApexDoorGivesRealWolfdorpHints covers round 159's real,
// sourced riddle content (see world.Room.DoorHints's doc comment) at
// Wolfdorp, an already-real, reachable CollodonsPile room.
func TestHandleApexDoorGivesRealWolfdorpHints(t *testing.T) {
	g := walkToWolfdorp(t)
	got := g.Handle(parser.Parse("APEX, DOOR"))
	for _, want := range []string{"Cry and enter door", "To enter is madness"} {
		if !strings.Contains(got, want) {
			t.Errorf("Handle(APEX, DOOR) at Wolfdorp = %q, want it to contain the real riddle %q", got, want)
		}
	}
}

func TestHandleSpeakToApexIsSynonymForTalk(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("APEX, SPEAK"))
	if !strings.Contains(got, "Apex") {
		t.Errorf("Handle(APEX, SPEAK) = %q, want the same recognition as APEX, TALK", got)
	}
}

// TestHandleApexThanksDismisses pins the hint screen's own confirmed
// dismiss phrase ("To dismiss say \"APEX, THANKS\"" — see game.help's
// verbatim text), wired for real for the first time this round.
func TestHandleApexThanksDismisses(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("APEX, THANKS"))
	if !strings.Contains(got, "thank Apex") {
		t.Errorf("Handle(APEX, THANKS) = %q, want a real dismiss response", got)
	}
}

// TestHandleCallWithoutScrollFails covers round 125's real, confirmed
// effect for CALL (see Handle's doc comment for the sourcing - The
// CRPG Addict's direct playthrough account): without the Scroll, it
// honestly fails rather than summoning Apex for free.
func TestHandleCallWithoutScrollFails(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("CALL"))
	if strings.Contains(got, "don't understand") {
		t.Errorf("Handle(CALL) = %q, want it recognized as a real spell, not an unknown word", got)
	}
	if strings.Contains(got, "Apex") {
		t.Errorf("Handle(CALL) without the Scroll = %q, want it to fail, not summon Apex", got)
	}
}

// TestHandleCallWithScrollSummonsApex covers CALL's real confirmed
// effect end-to-end: with the Scroll carried, it summons Apex.

// TestHandleCallWithScrollSummonsApex covers CALL's real confirmed
// effect end-to-end: with the Scroll carried, it summons Apex.
func TestHandleCallWithScrollSummonsApex(t *testing.T) {
	g := New()
	g.Player.Items = append(g.Player.Items, "Scroll")
	got := g.Handle(parser.Parse("CALL"))
	if !strings.Contains(got, "Apex") {
		t.Errorf("Handle(CALL) with the Scroll carried = %q, want it to summon Apex", got)
	}
}
