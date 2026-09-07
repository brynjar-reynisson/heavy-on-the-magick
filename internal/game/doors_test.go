package game

import (
	"strings"
	"testing"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/character"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
)

// TestPassGuardsClearsRealObstacle covers the confirmed real command
// "GUARDS, DOOR" (see world.Room.Guards and passGuards doc comments) —
// the same TARGET-comma-VERB grammar as "APEX, TALK".
func TestPassGuardsClearsRealObstacle(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Guards = true
	got := g.Handle(parser.Parse("GUARDS, DOOR"))
	if !strings.Contains(got, "let you pass") {
		t.Errorf("Handle(GUARDS, DOOR) with real guards present = %q, want them to let the player pass", got)
	}
	if g.World.CurrentRoom().Guards {
		t.Error("Guards should be cleared after a successful GUARDS, DOOR")
	}
}

func TestPassGuardsWithNoGuardsPresent(t *testing.T) {
	g := New() // Room of Misery has no Guards
	got := g.Handle(parser.Parse("GUARDS, DOOR"))
	if !strings.Contains(got, "no guards") {
		t.Errorf("Handle(GUARDS, DOOR) with no guards present = %q, want it to say there are none", got)
	}
}

func TestHandleDoorPasswordCorrect(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST")) // move to Secunda Porta, which has a door password
	got := g.Handle(parser.Parse("DOOR, SILENCE"))
	if !strings.Contains(got, "door swings open") {
		t.Errorf("Handle(DOOR, SILENCE) in Secunda Porta = %q, want it to succeed", got)
	}
}

func TestHandleDoorPasswordCaseInsensitive(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST"))
	got := g.Handle(parser.Parse("DOOR, silence"))
	if !strings.Contains(got, "door swings open") {
		t.Errorf("Handle(DOOR, silence) lowercase = %q, want it to still succeed", got)
	}
}

func TestHandleSecundaPortaDoorPromotesToZelator(t *testing.T) {
	g := New()
	if g.Player.Grade != character.Neophyte {
		t.Fatalf("test setup bug: new player Grade = %v, want Neophyte", g.Player.Grade)
	}
	g.Handle(parser.Parse("EAST")) // Secunda Porta
	got := g.Handle(parser.Parse("DOOR, SILENCE"))
	if g.Player.Grade != character.Zelator {
		t.Errorf("Grade after Secunda Porta's door = %v, want Zelator", g.Player.Grade)
	}
	if !strings.Contains(got, "Zelator") {
		t.Errorf("Handle(DOOR, SILENCE) = %q, want it to mention becoming a Zelator", got)
	}
}

func TestHandleDoorPasswordWrong(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST")) // Secunda Porta's real password is SILENCE, not THANKS
	got := g.Handle(parser.Parse("DOOR, THANKS"))
	if got != "Nothing happens." {
		t.Errorf("Handle(DOOR, THANKS) wrong password = %q, want rejection", got)
	}
}

func TestHandleDoorNoPasswordKnown(t *testing.T) {
	g := New() // Room of Misery has no known DoorPasswords
	got := g.Handle(parser.Parse("DOOR, ANYTHING"))
	if strings.Contains(got, "swings open") {
		t.Errorf("Handle(DOOR, ...) in a room with no known password = %q, should not succeed", got)
	}
}

// TestHandleTollDoorRequiresItem uses a synthetic TollItem room to
// isolate the "DOOR, <item>" form specifically (see
// TestHandleDropPaysRealToll below for the real, now-placed rooms and
// the "DROP <item>" form a fresh walkthrough re-read found to be the
// actual trigger phrase — see world.Room.TollItem's doc comment).
func TestHandleTollDoorRequiresItem(t *testing.T) {
	g := New()
	g.World.CurrentRoom().TollItem = "Bag of Gold"

	got := g.Handle(parser.Parse("DOOR, BAG OF GOLD"))
	if !strings.Contains(got, "don't have") {
		t.Errorf("Handle(DOOR, BAG OF GOLD) without the item = %q, want a don't-have response", got)
	}

	g.Player.Items = append(g.Player.Items, "Bag of Gold")
	got = g.Handle(parser.Parse("DOOR, BAG OF GOLD"))
	if !strings.Contains(got, "swings open") {
		t.Errorf("Handle(DOOR, BAG OF GOLD) carrying the item = %q, want it to succeed", got)
	}
	if g.hasItem("Bag of Gold") {
		t.Error("Bag of Gold should be spent (removed from inventory) after paying the toll")
	}
}

// TestHandleDropPaysRealToll covers the real, now-placed TollItem rooms
// (round 64) and confirms "DROP <item>" - the phrase a fresh CASA
// walkthrough re-read found is the actual trigger - opens the door.
func TestHandleDropPaysRealToll(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST")) // Secunda Porta
	g.Handle(parser.Parse("DOOR, SILENCE"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd
	g.Handle(parser.Parse("SOUTH")) // Sothic Complex
	g.Handle(parser.Parse("SOUTH")) // Wolfdorp
	room := g.World.CurrentRoom()
	if room.Name != "Wolfdorp" {
		t.Fatalf("test setup bug: expected to be in Wolfdorp, got %q", room.Name)
	}
	g.Handle(parser.Parse("PICKUP BAG"))
	g.Handle(parser.Parse("NORTH-WEST")) // Room of Stings
	if g.World.CurrentRoom().Name != "Room of Stings" {
		t.Fatalf("test setup bug: expected to be in Room of Stings, got %q", g.World.CurrentRoom().Name)
	}
	if got := g.Handle(parser.Parse("DROP BAG")); strings.Contains(got, "swings open") {
		t.Errorf("Handle(DROP BAG) at Room of Stings (needs a Key, not a Bag) = %q, should not open the door", got)
	}
	g.Handle(parser.Parse("PICKUP BAG")) // take it back before moving on

	g.Handle(parser.Parse("NORTH")) // Morfang, needs the Bag
	if g.World.CurrentRoom().Name != "Morfang" {
		t.Fatalf("test setup bug: expected to be in Morfang, got %q", g.World.CurrentRoom().Name)
	}
	got := g.Handle(parser.Parse("DROP BAG"))
	if !strings.Contains(got, "swings open") {
		t.Errorf("Handle(DROP BAG) at Morfang (needs a Bag) = %q, want it to open the door", got)
	}
	if g.hasItem("Bag") {
		t.Error("Bag should be spent (removed from inventory) after paying the real toll")
	}
}
