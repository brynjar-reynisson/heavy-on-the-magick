package world

import (
	"strings"
	"testing"
)

func TestRenderASCIIMapShowsOnlyVisitedRooms(t *testing.T) {
	w := PlaceholderDungeon()
	out := RenderASCIIMap(w)
	if out == "" {
		t.Fatal("RenderASCIIMap returned empty output for a world with a visited starting room")
	}

	w.Move(North) // visit "Narrow Passage"
	out = RenderASCIIMap(w)
	if len(out) == 0 {
		t.Fatal("RenderASCIIMap returned empty output after moving")
	}
}

func TestRenderASCIIMapMarksLivingMonster(t *testing.T) {
	w := CollodonsPile()
	w.Move(East)  // Secunda Porta
	w.Move(North) // Trollwynd, has a real living Monster
	out := RenderASCIIMap(w)
	if !strings.Contains(out, "[TRO]!") {
		t.Errorf("RenderASCIIMap with a living monster in Trollwynd = %q, want it marked with !", out)
	}
}

func TestRenderASCIIMapDoesNotMarkDefeatedMonster(t *testing.T) {
	w := CollodonsPile()
	w.Move(East)
	w.Move(North)                            // Trollwynd
	w.Rooms[roomTrollwynd].MonsterHealth = 0 // defeated
	out := RenderASCIIMap(w)
	if strings.Contains(out, "[TRO]!") {
		t.Errorf("RenderASCIIMap with a defeated monster = %q, should not still show !", out)
	}
}

func TestRenderASCIIMapMarksItems(t *testing.T) {
	w := CollodonsPile() // Room of Misery has a real item: Grimoire
	out := RenderASCIIMap(w)
	if !strings.Contains(out, "*") {
		t.Errorf("RenderASCIIMap with items in a visited room = %q, want an item marker (*)", out)
	}
}

// TestRenderASCIIMapMarksGuards covers the real Guards obstacle marker
// (see roomMarker's doc comment) using Level1Grid's confirmed D4 Guards
// placement (see level1_grid.go).
func TestRenderASCIIMapMarksGuards(t *testing.T) {
	w := Level1Grid()
	w.Move(South) // B1
	w.Move(South) // C1
	w.Move(East)  // C2
	w.Move(East)  // C3
	w.Move(East)  // C4
	w.Move(South) // D4, has a real Guards obstacle
	if !w.Rooms[w.Current].Guards {
		t.Fatalf("test setup bug: expected to be in a room with Guards, got %+v", w.Rooms[w.Current])
	}
	out := RenderASCIIMap(w)
	if !strings.Contains(out, "#") {
		t.Errorf("RenderASCIIMap with a real Guards obstacle = %q, want it marked with #", out)
	}
}

// TestRenderASCIIMapDoesNotMarkClearedGuards mirrors the defeated-
// monster test above for Guards: once cleared (game.passGuards sets
// Guards=false), the map should stop marking the room.
func TestRenderASCIIMapDoesNotMarkClearedGuards(t *testing.T) {
	w := Level1Grid()
	w.Move(South)
	w.Move(South)
	w.Move(East)
	w.Move(East)
	w.Move(East)
	w.Move(South) // D4
	w.Rooms[w.Current].Guards = false
	out := RenderASCIIMap(w)
	if strings.Contains(out, "#") {
		t.Errorf("RenderASCIIMap with cleared Guards = %q, should not still show #", out)
	}
}

// TestRenderASCIIMapMarksFire covers round 153's real Fire marker (see
// roomMarker's doc comment) - mirrors the already-shipped Guards test's
// synthetic-room pattern, since no reachable-from-a-fresh-game room
// currently carries Fire (Level2Grid's D6 is isolated).
func TestRenderASCIIMapMarksFire(t *testing.T) {
	w := New(0)
	w.AddRoom(&Room{ID: 0, Name: "Start", Exits: map[Direction]RoomID{North: 1}, Visited: true})
	w.AddRoom(&Room{ID: 1, Name: "Blaze", Fire: true})
	w.Move(North)

	out := RenderASCIIMap(w)
	if !strings.Contains(out, "[BLA]F") {
		t.Errorf("RenderASCIIMap with a real Fire hazard = %q, want it marked with F", out)
	}
}

// TestRenderASCIIMapMarksLockedDoor covers round 153's real locked-door
// marker (see roomMarker's doc comment), for both real gating fields
// (DoorPasswords and TollItem).
func TestRenderASCIIMapMarksLockedDoor(t *testing.T) {
	w := CollodonsPile() // Secunda Porta has a real DoorPasswords entry
	w.Move(East)

	out := RenderASCIIMap(w)
	if !strings.Contains(out, "[SEC]D") {
		t.Errorf("RenderASCIIMap with a real locked door = %q, want it marked with D", out)
	}
}

// TestRenderASCIIMapFireOutranksGuardsMarker covers roomMarker's stated
// priority order: Fire (which actually blocks movement) outranks Guards
// (which doesn't) when a room somehow has both.
func TestRenderASCIIMapFireOutranksGuardsMarker(t *testing.T) {
	w := New(0)
	w.AddRoom(&Room{ID: 0, Name: "Start", Exits: map[Direction]RoomID{North: 1}, Visited: true})
	w.AddRoom(&Room{ID: 1, Name: "Both", Fire: true, Guards: true})
	w.Move(North)

	if got := roomMarker(w.Rooms[w.Current]); got != "F" {
		t.Errorf("roomMarker with both Fire and Guards = %q, want F (Fire takes priority)", got)
	}
}

func TestRenderASCIIMapInconsistentFallsBackToList(t *testing.T) {
	w := New(0)
	a := &Room{ID: 0, Name: "A", Exits: map[Direction]RoomID{East: 1, SouthEast: 2}, Visited: true}
	b := &Room{ID: 1, Name: "B", Visited: true}
	c := &Room{ID: 2, Name: "C", Exits: map[Direction]RoomID{West: 1}, Visited: true}
	w.AddRoom(a)
	w.AddRoom(b)
	w.AddRoom(c)

	out := RenderASCIIMap(w)
	if out == "" {
		t.Fatal("expected fallback list output, got empty string")
	}
	// The fallback list format names rooms directly.
	for _, name := range []string{"A", "B", "C"} {
		found := false
		for _, r := range w.VisitedRooms() {
			if r.Name == name {
				found = true
			}
		}
		if !found {
			t.Fatalf("test setup bug: room %q not in VisitedRooms", name)
		}
	}
}
