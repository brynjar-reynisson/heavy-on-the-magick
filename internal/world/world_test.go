package world

import "testing"

func threeRoomLine() *World {
	// A -- East --> B -- East --> C, with matching West exits back.
	w := New(0)
	a := &Room{ID: 0, Name: "A", Exits: map[Direction]RoomID{East: 1}, Visited: true}
	b := &Room{ID: 1, Name: "B", Exits: map[Direction]RoomID{West: 0, East: 2}}
	c := &Room{ID: 2, Name: "C", Exits: map[Direction]RoomID{West: 1}}
	w.AddRoom(a)
	w.AddRoom(b)
	w.AddRoom(c)
	return w
}

func TestMove(t *testing.T) {
	w := threeRoomLine()

	if !w.Move(East) {
		t.Fatal("expected to move East from A to B")
	}
	if w.Current != 1 {
		t.Fatalf("Current = %d, want 1", w.Current)
	}
	if !w.Rooms[1].Visited {
		t.Error("room B should be marked Visited after moving into it")
	}

	if w.Move(North) {
		t.Error("expected Move(North) to fail: no such exit from B")
	}
}

func TestVisitedRooms(t *testing.T) {
	w := threeRoomLine()
	w.Move(East) // A -> B

	visited := w.VisitedRooms()
	if len(visited) != 2 {
		t.Fatalf("len(VisitedRooms()) = %d, want 2 (A and B, not C)", len(visited))
	}
}

func TestLayoutConsistent(t *testing.T) {
	w := threeRoomLine()
	w.Move(East) // visit B
	w.Move(East) // visit C

	positions, consistent := Layout(w, 0)
	if !consistent {
		t.Fatal("expected a simple East-West line to lay out consistently")
	}
	want := map[RoomID]Coord{0: {0, 0}, 1: {1, 0}, 2: {2, 0}}
	for id, wantPos := range want {
		if got := positions[id]; got != wantPos {
			t.Errorf("positions[%d] = %+v, want %+v", id, got, wantPos)
		}
	}
}

func TestLayoutInconsistent(t *testing.T) {
	// A -- East --> B, and also A -- East --> C, but C -- West --> B too:
	// B and C both claim the same grid cell East of A, so the layout can't
	// be trusted as a flat map (mirrors a non-Euclidean shortcut/loop in
	// the original dungeon design).
	w := New(0)
	a := &Room{ID: 0, Exits: map[Direction]RoomID{East: 1, SouthEast: 2}, Visited: true}
	b := &Room{ID: 1, Visited: true}
	c := &Room{ID: 2, Exits: map[Direction]RoomID{West: 1}, Visited: true}
	w.AddRoom(a)
	w.AddRoom(b)
	w.AddRoom(c)

	_, consistent := Layout(w, 0)
	if consistent {
		t.Fatal("expected layout to detect the conflicting position for room B")
	}
}
