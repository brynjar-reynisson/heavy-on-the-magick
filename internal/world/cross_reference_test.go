package world

import "testing"

// TestSharedNamedRoomsFiltersCellCodeCoincidences pins that plain
// row-letter+column-number names (e.g. "A1"), which every Level1-4Grid
// independently assigns to its own unnamed cells, are NOT reported as
// shared — they're a naming-scheme coincidence, not a real overlap. All
// 4 level grids have an "A1" cell; none of them are the same room.
func TestSharedNamedRoomsFiltersCellCodeCoincidences(t *testing.T) {
	shared := SharedNamedRooms(Level1Grid(), Level2Grid())
	for name := range shared {
		if cellCodePattern.MatchString(name) {
			t.Errorf("SharedNamedRooms reported cell-code coincidence %q as shared", name)
		}
	}
}

// TestSharedNamedRoomsCollodonsPileLevel1Grid pins the exact set of
// real, genuinely-named overlaps between CollodonsPile and Level1Grid
// this round's cross-world scan found — all 3 were already independently
// documented in Level1Grid's own doc comment before this test existed
// (this locks them in against an accidental future rename breaking the
// connection silently), and confirms the scan found no OTHERS beyond
// what's already been reasoned about.
func TestSharedNamedRoomsCollodonsPileLevel1Grid(t *testing.T) {
	shared := SharedNamedRooms(CollodonsPile(), Level1Grid())
	want := []string{"Agile Stair", "Room of Stings", "Room of Arrows"}
	if len(shared) != len(want) {
		t.Errorf("SharedNamedRooms(CollodonsPile, Level1Grid) = %v, want exactly %v", shared, want)
	}
	for _, name := range want {
		if _, ok := shared[name]; !ok {
			t.Errorf("SharedNamedRooms(CollodonsPile, Level1Grid) is missing %q", name)
		}
	}
}

// TestSharedNamedRoomsExitIsNotClaimedAsOneRoom documents, via a
// passing assertion rather than just a comment, that "Exit" appearing
// in both Level1Grid (G3) and Level4Grid (G2) is understood to be two
// DIFFERENT physical rooms (the dungeon has multiple real exits per
// round 19's win-condition find), not something SharedNamedRooms'
// caller should treat as a mergeable overlap the way Room of Stings/
// Room of Arrows/Agile Stair are.
func TestSharedNamedRoomsExitIsNotClaimedAsOneRoom(t *testing.T) {
	shared := SharedNamedRooms(Level1Grid(), Level4Grid())
	ids, ok := shared["Exit"]
	if !ok {
		t.Fatal(`SharedNamedRooms(Level1Grid, Level4Grid) is missing "Exit" - expected it to be found (just not treated as one room)`)
	}
	if ids[0] == ids[1] {
		t.Error(`Level1Grid's and Level4Grid's "Exit" rooms unexpectedly share a RoomID`)
	}
}
