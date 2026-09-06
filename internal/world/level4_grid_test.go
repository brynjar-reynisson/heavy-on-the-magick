package world

import "testing"

func TestLevel4GridHas22Cells(t *testing.T) {
	if len(level4Cells) != 22 {
		t.Fatalf("len(level4Cells) = %d, want 22 (the validated 17-cell main component plus 5 isolated named special rooms)", len(level4Cells))
	}
}

// TestLevel4GridTheChasmIsInMainComponent pins the round-58 fix: The
// Chasm (F4) was never actually isolated - it's a real cell in the
// main component that was just missing its name.
func TestLevel4GridTheChasmIsInMainComponent(t *testing.T) {
	w := Level4Grid()
	room := w.Rooms[level4Room("F4")]
	if room == nil || room.Name != "The Chasm" {
		t.Fatalf("room F4 = %+v, want Name \"The Chasm\"", room)
	}
	if len(room.Exits) == 0 {
		t.Error("The Chasm (F4) should have real Exits (part of the main connected component)")
	}
}

// TestLevel4GridIsolatedNamedRooms pins the 5 other real, tight-crop-
// verified named cells found in round 58 - genuinely absent from the
// main component, unlike The Chasm.
func TestLevel4GridIsolatedNamedRooms(t *testing.T) {
	w := Level4Grid()
	want := map[string]string{"D2": "Scales", "D3": "Doubt of Rabak", "F1": "The Crypt", "G2": "Exit", "G4": "Pride"}
	for code, name := range want {
		room := w.Rooms[level4Room(code)]
		if room == nil || room.Name != name {
			t.Errorf("room %s Name = %v, want %q", code, room, name)
		}
		if len(room.Exits) != 0 {
			t.Errorf("room %s Exits = %v, want none (connectivity not confirmed)", code, room.Exits)
		}
	}
}

func TestLevel4GridExitsAreReciprocal(t *testing.T) {
	w := Level4Grid()
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

// TestLevel4GridHasVerifiedMonster pins the one monster placement in
// this file - a Medusa at H5, tight-crop-verified against the map's own
// icon legend (see Level4Grid's doc comment).
func TestLevel4GridHasVerifiedMonster(t *testing.T) {
	w := Level4Grid()
	room := w.Rooms[level4Room("H5")]
	if room == nil || room.Monster != "Medusa" {
		t.Fatalf("room H5 Monster = %v, want \"Medusa\"", room)
	}
	if room.MonsterHealth <= 0 {
		t.Errorf("room H5 MonsterHealth = %d, want > 0", room.MonsterHealth)
	}
}

func TestLevel4GridIsFullyConnected(t *testing.T) {
	w := Level4Grid()
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
	if len(visited) != 17 {
		t.Errorf("only %d of 17 cells reachable from the start room, want all 17", len(visited))
	}
}
