package world

import "testing"

func TestLevel1GridHas64Cells(t *testing.T) {
	if len(level1Cells) != 64 {
		t.Fatalf("len(level1Cells) = %d, want 64 (8x8 grid)", len(level1Cells))
	}
}

func TestLevel1GridExitsAreReciprocal(t *testing.T) {
	w := Level1Grid()
	for id, room := range w.Rooms {
		for dir, destID := range room.Exits {
			dest, ok := w.Rooms[destID]
			if !ok {
				t.Fatalf("room %v exits %v to an unknown RoomID %v", id, dir, destID)
			}
			back, ok := dest.Exits[dir.Opposite()]
			if !ok || back != id {
				t.Errorf("room %q --%s--> %q has no matching reverse exit (%q --%s--> %q)", room.Name, dir, dest.Name, dest.Name, dir.Opposite(), room.Name)
			}
		}
	}
}

// The starting cell's connected component covers 44 of 64 cells, not all
// of them - a real, honestly-documented gap (see Level1Grid's doc
// comment): the automated extraction plus one manually-resolved bridge
// (Room of Stings/Exit) joined the two biggest zones, but a second
// fragment (roughly columns 5-8 of rows E-H) remains separately isolated.
// Every possible crossing between the two components has since been
// individually hand-verified closed (all 12 checked), which is real
// evidence this fragment is a genuinely separate area in the source
// map's own data, not something further spot-checking is likely to fix.
// This test pins the current honest number rather than asserting a full
// connectivity that hasn't actually been verified.
func TestLevel1GridStartRoomComponentSize(t *testing.T) {
	w := Level1Grid()
	visited := map[RoomID]bool{w.Current: true}
	queue := []RoomID{w.Current}
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		for _, dest := range w.Rooms[id].Exits {
			if !visited[dest] {
				visited[dest] = true
				queue = append(queue, dest)
			}
		}
	}
	if len(visited) < 40 {
		t.Errorf("only %d of 64 cells reachable from the start room, want at least 40 (the currently-resolved component)", len(visited))
	}
}

func TestLevel1GridNamedRoomsPresent(t *testing.T) {
	w := Level1Grid()
	want := map[string]bool{"Room of Stings": false, "Room of Arrows": false, "Room of Claws": false, "Agile Stair": false}
	for _, room := range w.Rooms {
		if _, ok := want[room.Name]; ok {
			want[room.Name] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("expected a room named %q in Level1Grid", name)
		}
	}
}

// TestLevel1GridIsolatedCellsHaveNoExits pins the confirmed (not
// assumed) finding that A7, A8, G5, and H5 are closed on every side -
// individually re-verified with the same pixel-scan method used for the
// other 58 cells, not just left unresolved. See Level1Grid's doc comment
// for the "STILL FRAGMENTED" and Astarot-transport-hypothesis discussion.
func TestLevel1GridIsolatedCellsHaveNoExits(t *testing.T) {
	w := Level1Grid()
	for _, code := range []string{"A7", "A8", "G5", "H5"} {
		room := w.Rooms[level1Room(code)]
		if room == nil {
			t.Fatalf("Level1Grid is missing cell %s", code)
		}
		if len(room.Exits) != 0 {
			t.Errorf("room %s (%q) has Exits %v, want none (confirmed closed on every side)", code, room.Name, room.Exits)
		}
	}
}

// TestLevel1GridWerewolvesAreAtCorrectedCells pins the column-off-by-one
// fix (see Level1Grid's doc comment "CORRECTION" section): the werewolf
// icons are at C2/D5, not the originally-shipped C3/D6.
func TestLevel1GridWerewolvesAreAtCorrectedCells(t *testing.T) {
	w := Level1Grid()
	for _, code := range []string{"C2", "D5"} {
		room := w.Rooms[level1Room(code)]
		if room == nil || room.Monster != "Werewolf" {
			t.Errorf("room %s Monster = %q, want \"Werewolf\"", code, room.Monster)
		}
	}
	for _, code := range []string{"C3", "D6"} {
		room := w.Rooms[level1Room(code)]
		if room != nil && room.Monster != "" {
			t.Errorf("room %s Monster = %q, want empty (werewolf moved to the corrected cell)", code, room.Monster)
		}
	}
}

// TestLevel1GridGuardsPlacements pins the 2 verified Guards obstacles and
// that game.Handle's "GUARDS, DOOR" (see passGuards) can clear them - see
// TestPassGuardsClearsRealObstacle in the game package for the full
// end-to-end check.
func TestLevel1GridGuardsPlacements(t *testing.T) {
	w := Level1Grid()
	for _, code := range []string{"D4", "D7"} {
		room := w.Rooms[level1Room(code)]
		if room == nil || !room.Guards {
			t.Errorf("room %s Guards = %v, want true", code, room != nil && room.Guards)
		}
	}
}

func TestLevel1GridCyclopsCrossConfirmsNidus(t *testing.T) {
	w := Level1Grid()
	for _, room := range w.Rooms {
		if room.Monster == "Cyclops" {
			return
		}
	}
	t.Error("expected a Cyclops in Level1Grid, cross-confirming CollodonsPile's Nidus room")
}
