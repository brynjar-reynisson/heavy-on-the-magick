package world

import "testing"

func TestLevel4GridHas27Cells(t *testing.T) {
	if len(level4Cells) != 27 {
		t.Fatalf("len(level4Cells) = %d, want 27 (the validated 17-cell main component plus 10 isolated cells)", len(level4Cells))
	}
}

// TestLevel4GridWormringWyverns pins the round-72 find: 3 more real,
// tight-crop-verified Wyverns at A1/A3/A5, which together with the
// already-shipped E5 exactly match zone_monsters.go's "Wormring: Wyvern
// x4" sighting - real cross-validation between this map's icon data and
// that independent source. Also pins the real Vampire found at A6
// (Methos zone), corroborating the prior round's CollodonsPile Methos
// Vampire placement from a completely different source. ("Vampire" is
// the round-74 rename of what this project called "Wraith" - see
// Level1Grid's doc comment.)
func TestLevel4GridWormringWyverns(t *testing.T) {
	w := Level4Grid()
	wyverns := []string{"A1", "A3", "A5", "E5"}
	for _, code := range wyverns {
		room := w.Rooms[level4Room(code)]
		if room == nil || room.Monster != "Wyvern" || room.MonsterHealth <= 0 {
			t.Errorf("room %s = %+v, want an isolated cell with a live Wyvern", code, room)
		}
	}
	a6 := w.Rooms[level4Room("A6")]
	if a6 == nil || a6.Monster != "Vampire" || a6.MonsterHealth <= 0 {
		t.Errorf("room A6 = %+v, want an isolated cell with a live Vampire", a6)
	}
}

// TestLevel4GridHasSecondWyvern pins the round-68 find: a real,
// tight-crop-verified Wyvern at E5, isolated (no Exits) same as the
// round-58 finds.
func TestLevel4GridHasSecondWyvern(t *testing.T) {
	w := Level4Grid()
	room := w.Rooms[level4Room("E5")]
	if room == nil || room.Monster != "Wyvern" || room.MonsterHealth <= 0 {
		t.Fatalf("room E5 = %+v, want an isolated cell with a live Wyvern", room)
	}
	if len(room.Exits) != 0 {
		t.Errorf("room E5 Exits = %v, want none (connectivity not confirmed)", room.Exits)
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
	// Round 79: zone_monsters.go's "Doubt of Rabak: Vampire x1" sighting,
	// a separate entry from Methos's (already corroborated via A6).
	d3 := w.Rooms[level4Room("D3")]
	if d3.Monster != "Vampire" || d3.MonsterHealth <= 0 {
		t.Errorf("Doubt of Rabak (D3) Monster = %q (health %d), want a live Vampire", d3.Monster, d3.MonsterHealth)
	}
	// Round 183: user-recalled directly from finishing the same third
	// gameplay video independently - "Rabak goes down when we say
	// water" and "is impossible to pass until the correct words are
	// spoken" - see world.Room.Chasm's doc comment for the sourcing
	// discipline used here.
	if !d3.Water {
		t.Error("Doubt of Rabak (D3) should have a real Water hazard")
	}
}

// TestLevel4GridTheChasmIsLethalWithoutFlask pins round 183's real,
// user-recalled hazard: The Chasm (F4) is now a real, sourced deadly
// obstacle without a Flask - see world.Room.Chasm's doc comment.
func TestLevel4GridTheChasmIsLethalWithoutFlask(t *testing.T) {
	w := Level4Grid()
	room := w.Rooms[level4Room("F4")]
	if !room.Chasm {
		t.Error("The Chasm (F4) should have Chasm=true")
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
