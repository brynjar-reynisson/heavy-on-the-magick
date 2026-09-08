package game

import (
	"strings"
	"testing"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/character"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/world"
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

// TestPassWaterClearsRealObstacle covers round 169's real, sourced
// command "WATER, FALL" (CRASH 31's Signpost column: "To get past the
// water say 'Water, fall'") - see world.Room.Water and passWater's doc
// comments.
func TestPassWaterClearsRealObstacle(t *testing.T) {
	g := New()
	g.World.CurrentRoom().Water = true
	got := g.Handle(parser.Parse("WATER, FALL"))
	if !strings.Contains(got, "Trickle") {
		t.Errorf("Handle(WATER, FALL) with real water present = %q, want the real confirmed response \"Trickle\"", got)
	}
	if g.World.CurrentRoom().Water {
		t.Error("Water should be cleared after a successful WATER, FALL")
	}
}

func TestPassWaterWithNoWaterPresent(t *testing.T) {
	g := New() // Room of Misery has no Water
	got := g.Handle(parser.Parse("WATER, FALL"))
	if !strings.Contains(got, "no water") {
		t.Errorf("Handle(WATER, FALL) with no water present = %q, want it to say there is none", got)
	}
}

// TestLevel3GridWaterHazard covers the real, exact-cell placement (H4,
// already independently confirmed as literally named "Water") end to
// end via a direct Teleport - H4 is currently isolated (no Exits), the
// same honest "real but not live-walkthrough-reachable" scope several
// other isolated named cells in this project have.
func TestLevel3GridWaterHazard(t *testing.T) {
	g := NewLevel3Exploration()
	id, ok := g.World.FindRoomByName("Water")
	if !ok {
		t.Fatal("test setup bug: Level3Grid has no room named Water")
	}
	g.World.Teleport(id)
	if !g.World.CurrentRoom().Water {
		t.Fatalf("test setup bug: expected a real Water hazard at %q, got %+v", "Water", g.World.CurrentRoom())
	}
	got := g.Handle(parser.Parse("WATER, FALL"))
	if !strings.Contains(got, "Trickle") {
		t.Errorf("Handle(WATER, FALL) at the real Water cell = %q, want the real confirmed response \"Trickle\"", got)
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

// TestHandleTertiaPortaDoorPromotesToPracticus covers round 178's real,
// video-confirmed second promotion door: passing Tertia Porta's door
// raises Axil from Zelator to Practicus, the exact same pattern as
// Secunda Porta's Neophyte-to-Zelator promotion. No shipped World.Room
// is named "Tertia Porta" yet (a known real room name per
// known_room_names.go, not placed at an exact cell in any dataset) -
// uses a synthetic room, the same "mechanic real, not yet reachable"
// pattern as several other confirmed mechanics in this project.
func TestHandleTertiaPortaDoorPromotesToPracticus(t *testing.T) {
	w := world.New(0)
	w.AddRoom(&world.Room{ID: 0, Name: "Tertia Porta", DoorPasswords: []string{"OPEN"}})
	g := &Game{Player: character.NewPlayer(), World: w}
	g.Player.Grade = character.Zelator

	got := g.Handle(parser.Parse("DOOR, OPEN"))
	if g.Player.Grade != character.Practicus {
		t.Errorf("Grade after Tertia Porta's door = %v, want Practicus", g.Player.Grade)
	}
	if !strings.Contains(got, "Practicus") {
		t.Errorf("Handle(DOOR, OPEN) at Tertia Porta = %q, want it to mention becoming a Practicus", got)
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
