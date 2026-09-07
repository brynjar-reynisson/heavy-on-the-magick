package game

import (
	"strings"
	"testing"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
)

func TestHandleUnknownWord(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("GOLANG, COMPILE"))
	if !strings.Contains(got, "don't understand") {
		t.Errorf("Handle(unknown word) = %q, want a not-understood response", got)
	}
}

func TestHandleKnownWordNoBehaviorYet(t *testing.T) {
	g := New()
	// TALK to MONSTER (as opposed to APEX, which now has real behavior -
	// see TestHandleTalkToApex) is real confirmed vocabulary but has no
	// verb resolution yet.
	got := g.Handle(parser.Parse("MONSTER, TALK"))
	if !strings.Contains(got, "don't know what it does") {
		t.Errorf("Handle(known word, unimplemented verb) = %q, want the not-yet-implemented response", got)
	}
}

// walkToWolfdorp navigates a fresh Game from Room of Misery to Wolfdorp
// via real moves - the shared setup TestHandleDropPaysRealToll already
// uses, split out so other Wolfdorp-specific tests (like the HasChest
// ones below) don't have to repeat it.
func walkToWolfdorp(t *testing.T) *Game {
	t.Helper()
	g := New()
	g.Handle(parser.Parse("EAST"))          // Secunda Porta
	g.Handle(parser.Parse("DOOR, SILENCE")) // unlocks the door North
	g.Handle(parser.Parse("NORTH"))         // Trollwynd
	g.Handle(parser.Parse("SOUTH"))         // Sothic Complex
	g.Handle(parser.Parse("SOUTH"))         // Wolfdorp
	if room := g.World.CurrentRoom(); room == nil || room.Name != "Wolfdorp" {
		t.Fatalf("test setup bug: expected to be in Wolfdorp, got %+v", room)
	}
	return g
}

// TestHandleLookMentionsChest covers the real HasChest fixture (round
// 78, see world.Room.HasChest's doc comment) being surfaced in LOOK
// itself, not just discoverable by blindly guessing "EXAMINE CHEST".

func TestLevel1ExplorationStartRoomIsVisited(t *testing.T) {
	// Regression test: NewLevel1Exploration originally forgot to mark its
	// start room Visited (unlike New()/CollodonsPile), which crashed
	// Handle(MAP) with a nil pointer dereference the moment it was
	// actually run (caught by playing it live, not by unit tests alone -
	// see CLAUDE.md).
	g := NewLevel1Exploration()
	if !g.World.CurrentRoom().Visited {
		t.Fatal("NewLevel1Exploration's start room should be marked Visited immediately")
	}
	got := g.Handle(parser.Parse("MAP"))
	if got == "" {
		t.Error("Handle(MAP) on a fresh Level1Exploration game returned an empty string")
	}
}

func TestLevel2ExplorationStartRoomIsVisited(t *testing.T) {
	g := NewLevel2Exploration()
	if !g.World.CurrentRoom().Visited {
		t.Fatal("NewLevel2Exploration's start room should be marked Visited immediately")
	}
	got := g.Handle(parser.Parse("MAP"))
	if got == "" {
		t.Error("Handle(MAP) on a fresh Level2Exploration game returned an empty string")
	}
}

func TestLevel3ExplorationStartRoomIsVisited(t *testing.T) {
	g := NewLevel3Exploration()
	if !g.World.CurrentRoom().Visited {
		t.Fatal("NewLevel3Exploration's start room should be marked Visited immediately")
	}
	got := g.Handle(parser.Parse("MAP"))
	if got == "" {
		t.Error("Handle(MAP) on a fresh Level3Exploration game returned an empty string")
	}
}

func TestLevel4ExplorationStartRoomIsVisited(t *testing.T) {
	g := NewLevel4Exploration()
	if !g.World.CurrentRoom().Visited {
		t.Fatal("NewLevel4Exploration's start room should be marked Visited immediately")
	}
	got := g.Handle(parser.Parse("MAP"))
	if got == "" {
		t.Error("Handle(MAP) on a fresh Level4Exploration game returned an empty string")
	}
}

func TestNewGameStartsInRoomOfMisery(t *testing.T) {
	g := New()
	room := g.World.CurrentRoom()
	if room == nil || room.Name != "Room of Misery" {
		t.Errorf("starting room = %+v, want Room of Misery", room)
	}
	if !room.Visited {
		t.Error("starting room should be marked Visited immediately")
	}
}
