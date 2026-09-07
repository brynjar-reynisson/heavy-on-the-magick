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

// TestRenderASCIIMapMarksWater covers round 170's real Water marker
// (see roomMarker's doc comment) using Level3Grid's confirmed H4 Water
// placement (round 169) - reached via a direct Teleport since H4 is
// currently isolated (no Exits), the same honest scope every other
// isolated named cell in this project has had before its own
// connectivity is found. Deliberately checks for the exact "> W Water"
// marker position (the real list-view format, since round 170's own
// fix routes a Teleport-only-reached isolated room through the list
// fallback - see RenderASCIIMap's doc comment), not just a bare "W" -
// the room's own full name ("Water") would make a bare "W" check pass
// even with no real marker at all, a false-positive risk caught while
// writing this exact test.
func TestRenderASCIIMapMarksWater(t *testing.T) {
	w := Level3Grid()
	id, ok := w.FindRoomByName("Water")
	if !ok {
		t.Fatal("test setup bug: Level3Grid has no room named Water")
	}
	w.Teleport(id)
	if !w.Rooms[w.Current].Water {
		t.Fatalf("test setup bug: expected to be in a room with Water, got %+v", w.Rooms[w.Current])
	}
	out := RenderASCIIMap(w)
	if !strings.Contains(out, "> W Water") {
		t.Errorf("RenderASCIIMap with a real Water hazard = %q, want it marked with W", out)
	}
}

// TestRenderASCIIMapDoesNotMarkClearedWater mirrors the cleared-Guards
// test above for Water: once cleared (game.passWater sets Water=false),
// the map should stop marking the room. Checks for the ABSENCE of the
// exact marked pattern, not a bare "W" (see TestRenderASCIIMapMarksWater's
// doc comment for why a bare check would be a false-positive risk here).
func TestRenderASCIIMapDoesNotMarkClearedWater(t *testing.T) {
	w := Level3Grid()
	id, ok := w.FindRoomByName("Water")
	if !ok {
		t.Fatal("test setup bug: Level3Grid has no room named Water")
	}
	w.Teleport(id)
	w.Rooms[w.Current].Water = false
	out := RenderASCIIMap(w)
	if strings.Contains(out, "> W Water") {
		t.Errorf("RenderASCIIMap with cleared Water = %q, should not still show the W marker", out)
	}
}

// TestRenderASCIIMapShowsTeleportOnlyVisitedRoom covers round 170's
// real bug fix: Layout's BFS only reaches rooms connected via real
// Exits from the anchor - a genuinely visited room reached ONLY via
// World.Teleport (no Exits in or out) used to be silently omitted from
// the map entirely, with `consistent` staying true the whole time (no
// contradiction was ever detected, since Layout's BFS simply never
// visited that room to find one). Confirms the real fix: the room now
// shows up (via the list-view fallback), not silently vanishes.
func TestRenderASCIIMapShowsTeleportOnlyVisitedRoom(t *testing.T) {
	w := Level3Grid()
	id, ok := w.FindRoomByName("Water")
	if !ok {
		t.Fatal("test setup bug: Level3Grid has no room named Water")
	}
	w.Teleport(id)
	out := RenderASCIIMap(w)
	if !strings.Contains(out, "Water") {
		t.Errorf("RenderASCIIMap after teleporting to an Exit-less room = %q, want the room to still appear, not silently vanish", out)
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
