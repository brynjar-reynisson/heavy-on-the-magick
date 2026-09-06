package world

import "testing"

func TestLevel2GridHas52Cells(t *testing.T) {
	if len(level2Cells) != 52 {
		t.Fatalf("len(level2Cells) = %d, want 52 (the validated 50-cell main component plus 2 isolated named cells from the Room of Misery pocket)", len(level2Cells))
	}
}

// TestLevel2GridRoomOfMiseryPocketNamedCells pins the 2 real, tight-crop-
// verified named cells added from Room of Misery's pocket (see this
// file's doc comment) - F4 is the game's confirmed real starting room.
func TestLevel2GridRoomOfMiseryPocketNamedCells(t *testing.T) {
	w := Level2Grid()
	want := map[string]string{"F3": "Sign", "F4": "Room of Misery"}
	for code, name := range want {
		room := w.Rooms[level2Room(code)]
		if room == nil || room.Name != name {
			t.Errorf("room %s Name = %v, want %q", code, room, name)
		}
		if len(room.Exits) != 0 {
			t.Errorf("room %s Exits = %v, want none (connectivity not confirmed)", code, room.Exits)
		}
	}
}

func TestLevel2GridExitsAreReciprocal(t *testing.T) {
	w := Level2Grid()
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

// TestLevel2GridHasVerifiedMonsters pins the 3 monster placements added
// after Level2Grid's initial (connectivity-only) shipment — each
// individually verified with a tight per-cell crop against the icon
// legend, not eyeballed from the full grid view.
func TestLevel2GridHasVerifiedMonsters(t *testing.T) {
	w := Level2Grid()
	want := map[string]string{"A5": "Wraith", "C2": "Slug", "H6": "Ghost"}
	for code, monster := range want {
		room := w.Rooms[level2Room(code)]
		if room == nil {
			t.Fatalf("Level2Grid is missing cell %s", code)
		}
		if room.Monster != monster {
			t.Errorf("room %s Monster = %q, want %q", code, room.Monster, monster)
		}
		if room.MonsterHealth <= 0 {
			t.Errorf("room %s MonsterHealth = %d, want > 0", code, room.MonsterHealth)
		}
	}
}

// TestLevel2GridGuardsPlacements pins the 3 verified Guards obstacles
// within the playable 50-cell component - see Level2Grid's doc comment
// for why 2 more real Guards icons (G5, H5) are deliberately excluded.
func TestLevel2GridGuardsPlacements(t *testing.T) {
	w := Level2Grid()
	for _, code := range []string{"B1", "C8", "D8"} {
		room := w.Rooms[level2Room(code)]
		if room == nil || !room.Guards {
			t.Errorf("room %s Guards = %v, want true", code, room != nil && room.Guards)
		}
	}
}

// TestLevel2GridIsFullyConnected pins that, unlike Level1Grid, this
// 50-cell component is entirely reachable from the start room - it was
// deliberately built from one whole connected component (see
// Level2Grid's doc comment), not assembled from a partially-connected
// automated pass.
func TestLevel2GridIsFullyConnected(t *testing.T) {
	w := Level2Grid()
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
	if len(visited) != 50 {
		t.Errorf("only %d of 50 cells reachable from the start room, want all 50", len(visited))
	}
}
