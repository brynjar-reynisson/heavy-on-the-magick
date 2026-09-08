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
// this project's cross-world scan found — the first 3 were already
// independently documented in Level1Grid's own doc comment before this
// test existed; "Furnace Room" joined in round 126 (CollodonsPile's
// own doc comment - a real, independently-sourced match to Level1Grid's
// A8, cross-confirmed on BOTH name and "no exits"). "Exit" joined once
// CollodonsPile got its own real Exit room (Pile Collodom's North exit,
// confirmed via a full walkthrough video) - same honest treatment as
// TestSharedNamedRoomsExitIsNotClaimedAsOneRoom below: a real name
// match, not proof these are literally the same physical room (the
// dungeon has multiple real exits). Locks all 5 in against an
// accidental future rename breaking the connection silently, and
// confirms the scan finds no OTHERS beyond what's already reasoned
// about.
func TestSharedNamedRoomsCollodonsPileLevel1Grid(t *testing.T) {
	shared := SharedNamedRooms(CollodonsPile(), Level1Grid())
	want := []string{"Agile Stair", "Room of Stings", "Room of Arrows", "Furnace Room", "Exit"}
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
