package world

import "testing"

func TestLevel2GridHas61Cells(t *testing.T) {
	if len(level2Cells) != 61 {
		t.Fatalf("len(level2Cells) = %d, want 61 (the validated 50-cell main component, the full 7-cell Room of Misery pocket, 3 isolated named cells: Icthys/Horns/Purity, and D6's Fire hazard)", len(level2Cells))
	}
}

// TestLevel2GridFloxIsInMainComponent pins the round-56 fix: Flox (D4)
// was never actually isolated - it's a real cell in the main component
// that was just missing its name (see this file's doc comment).
func TestLevel2GridFloxIsInMainComponent(t *testing.T) {
	w := Level2Grid()
	room := w.Rooms[level2Room("D4")]
	if room == nil || room.Name != "Flox" {
		t.Fatalf("room D4 = %+v, want Name \"Flox\"", room)
	}
	if len(room.Exits) == 0 {
		t.Error("Flox (D4) should have real Exits (it's part of the main connected component, not isolated)")
	}
}

// TestLevel2GridA1IsExit pins round 130's find: A1 - this file's own
// arbitrary starting anchor - is a real, tight-crop-verified "Exit"
// (see this file's doc comment), the 3rd of the dungeon's confirmed
// 3 real Exit rooms this project now ships (alongside Level1Grid's G3
// and Level4Grid's G2).
func TestLevel2GridA1IsExit(t *testing.T) {
	w := Level2Grid()
	room := w.Rooms[level2Room("A1")]
	if room == nil || room.Name != "Exit" {
		t.Fatalf("room A1 = %+v, want Name \"Exit\"", room)
	}
	if len(room.Exits) != 2 {
		t.Errorf("Exit (A1) Exits = %v, want its real East/South connectivity unchanged", room.Exits)
	}
}

// TestLevel2GridD6HasFire pins the round-80 find: a real, previously
// entirely-missing cell (D6) confirmed via a tight-crop-verified
// "FIRE!" warning label (see world.Room.Fire's doc comment).
func TestLevel2GridD6HasFire(t *testing.T) {
	w := Level2Grid()
	room := w.Rooms[level2Room("D6")]
	if room == nil || !room.Fire {
		t.Errorf("room D6 = %+v, want Fire = true", room)
	}
	if len(room.Exits) != 0 {
		t.Errorf("room D6 Exits = %v, want none (connectivity not confirmed)", room.Exits)
	}
}

// TestLevel2GridOtherNamedIsolatedCells pins Icthys/Horns/Purity - real,
// tight-crop-verified named cells that (unlike Flox) really are absent
// from the main component.
func TestLevel2GridOtherNamedIsolatedCells(t *testing.T) {
	w := Level2Grid()
	want := map[string]string{"C3": "Icthys", "E2": "Horns", "H3": "Purity"}
	for code, name := range want {
		room := w.Rooms[level2Room(code)]
		if room == nil || room.Name != name {
			t.Errorf("room %s Name = %v, want %q", code, room, name)
		}
		if len(room.Exits) != 0 {
			t.Errorf("room %s Exits = %v, want none (connectivity not confirmed)", code, room.Exits)
		}
	}
	// Round 71: zone_monsters.go's "Room of Icthys: Slug x1" sighting,
	// cross-validated against other exact matches in that same list.
	icthys := w.Rooms[level2Room("C3")]
	if icthys.Monster != "Slug" || icthys.MonsterHealth <= 0 {
		t.Errorf("Icthys (C3) Monster = %q (health %d), want a live Slug", icthys.Monster, icthys.MonsterHealth)
	}
}

// TestLevel2GridRoomOfMiseryPocketNamedCells pins the 2 real, tight-crop-
// verified named cells in Room of Misery's pocket (see this file's doc
// comment) - F4 is the game's confirmed real starting room.
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

// TestLevel2GridRoomOfMiseryPocketIsComplete pins that all 7 pocket
// cells now exist (round 55 added the remaining G3/G4/G5/H4/H5), with
// real, tight-crop-verified Guards obstacles at G5 and H5.
func TestLevel2GridRoomOfMiseryPocketIsComplete(t *testing.T) {
	w := Level2Grid()
	for _, code := range []string{"G3", "G4", "G5", "H4", "H5"} {
		if w.Rooms[level2Room(code)] == nil {
			t.Errorf("Level2Grid is missing pocket cell %s", code)
		}
	}
	for _, code := range []string{"G5", "H5"} {
		room := w.Rooms[level2Room(code)]
		if room == nil || !room.Guards {
			t.Errorf("room %s Guards = %v, want true", code, room != nil && room.Guards)
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
	want := map[string]string{"A5": "Vampire", "C2": "Slug", "H6": "Ghost"}
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

// TestLevel2GridA5HasSnake pins round 179's addition: a full frame-by-
// frame review of real gameplay footage shows the player picking up a
// real Snake at "Wraithvale" (A5's own real zone name, round 74) - the
// same real, sourced ward-off item game.checkSnakeHydra already
// requires to pass a Hydra, previously shipped with no real placement
// anywhere.
func TestLevel2GridA5HasSnake(t *testing.T) {
	w := Level2Grid()
	room := w.Rooms[level2Room("A5")]
	found := false
	for _, item := range room.Items {
		if item == "Snake" {
			found = true
		}
	}
	if !found {
		t.Errorf("Level2Grid room A5 Items = %v, want it to include \"Snake\"", room.Items)
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
