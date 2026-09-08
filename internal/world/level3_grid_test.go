package world

import "testing"

func TestLevel3GridHas48Cells(t *testing.T) {
	if len(level3Cells) != 48 {
		t.Fatalf("len(level3Cells) = %d, want 48 (the validated 41-cell connected component plus 7 isolated cells: Sothic Complex, Room of Nani, Hydra, Two, G4/Wyvern, Kitchen of Ai, Water)", len(level3Cells))
	}
}

// TestLevel3GridKitchenOfAiFinds pins the 3 real, tight-crop-verified
// finds from round 66 (see Level3Grid's doc comment): Two (G2), a
// Wyvern monster (G4), and Water (H4) - all isolated, same convention
// as Sothic Complex/Nani/Hydra.
func TestLevel3GridKitchenOfAiFinds(t *testing.T) {
	w := Level3Grid()
	two := w.Rooms[level3Room("G2")]
	if two == nil || two.Name != "Two" || len(two.Exits) != 0 {
		t.Errorf("room G2 = %+v, want isolated with Name \"Two\"", two)
	}
	g4 := w.Rooms[level3Room("G4")]
	if g4 == nil || g4.Monster != "Wyvern" || g4.MonsterHealth <= 0 || len(g4.Exits) != 0 {
		t.Errorf("room G4 = %+v, want isolated with a live Wyvern", g4)
	}
	water := w.Rooms[level3Room("H4")]
	if water == nil || water.Name != "Water" || len(water.Exits) != 0 {
		t.Errorf("room H4 = %+v, want isolated with Name \"Water\"", water)
	}
}

// TestLevel3GridKitchenOfAiRoomHasCauldron pins the round-178 addition:
// a full frame-by-frame review of real gameplay footage shows a real,
// distinct room "Kitchen of Ai" (not just the zone label round 66
// already knew about) with a real, examinable Cauldron holding a
// Scroll - the real location for game.cauldronAchad's "CAULDRON,
// ACHAD" ritual (corrected from round 177's "Room of Nani" guess).
func TestLevel3GridKitchenOfAiRoomHasCauldron(t *testing.T) {
	w := Level3Grid()
	room := w.Rooms[level3Room("H2")]
	if room == nil || room.Name != "Kitchen of Ai" || !room.HasCauldron || len(room.Exits) != 0 {
		t.Errorf("room H2 = %+v, want isolated, named \"Kitchen of Ai\", with a real Cauldron", room)
	}
	found := false
	for _, item := range room.Items {
		if item == "Scroll" {
			found = true
		}
	}
	if !found {
		t.Errorf("Kitchen of Ai Items = %v, want it to include \"Scroll\"", room.Items)
	}
}

// TestLevel3GridNaniAndHydraAreIsolated pins the 2 real, tight-crop-
// verified named cells found in round 57 (see Level3Grid's doc
// comment) - the same "check for a named special room" technique that
// found Sothic Complex, applied to the rest of the F2-F6 gap. F3's
// name was corrected round 177 from "Nani" to its real, full,
// video-confirmed status-panel name "Room of Nani".
func TestLevel3GridNaniAndHydraAreIsolated(t *testing.T) {
	w := Level3Grid()
	want := map[string]string{"F3": "Room of Nani", "F5": "Hydra"}
	for code, name := range want {
		room := w.Rooms[level3Room(code)]
		if room == nil || room.Name != name {
			t.Errorf("room %s Name = %v, want %q", code, room, name)
		}
		if len(room.Exits) != 0 {
			t.Errorf("room %s Exits = %v, want none (connectivity not confirmed)", code, room.Exits)
		}
	}
	// Round 71: zone_monsters.go's "Rook of Hydra: Wyvern x1" sighting,
	// cross-validated against other exact matches in that same list.
	hydra := w.Rooms[level3Room("F5")]
	if hydra.Monster != "Wyvern" || hydra.MonsterHealth <= 0 {
		t.Errorf("Hydra (F5) Monster = %q (health %d), want a live Wyvern", hydra.Monster, hydra.MonsterHealth)
	}
}

// TestLevel3GridSothicComplexIsIsolated pins the real, named-but-
// unconnected D4 cell (see Level3Grid's doc comment: it's drawn as an
// irregular special room, not a standard grid box, which is why the
// automated border-detection pass never found exits for it - the same
// situation as Level1Grid's A7/A8/G5/H5).
func TestLevel3GridSothicComplexIsIsolated(t *testing.T) {
	w := Level3Grid()
	room := w.Rooms[level3Room("D4")]
	if room == nil || room.Name != "Sothic Complex" {
		t.Fatalf("expected a room named Sothic Complex at D4, got %+v", room)
	}
	if len(room.Exits) != 0 {
		t.Errorf("Sothic Complex (D4) has Exits %v, want none (connectivity not confirmed)", room.Exits)
	}
	// Round 71: zone_monsters.go's "Sothic Complex: Ghost x1" sighting,
	// cross-validated against other exact matches in that same list.
	if room.Monster != "Ghost" || room.MonsterHealth <= 0 {
		t.Errorf("Sothic Complex (D4) Monster = %q (health %d), want a live Ghost", room.Monster, room.MonsterHealth)
	}
}

func TestLevel3GridExitsAreReciprocal(t *testing.T) {
	w := Level3Grid()
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

// TestLevel3GridHasVerifiedMonsters pins the 7 monster placements added
// after the calibration fix — each individually tight-crop-verified
// against the map's own icon legend (see Level3Grid's doc comment for
// the color-ambiguity/false-positive details this caught).
func TestLevel3GridHasVerifiedMonsters(t *testing.T) {
	w := Level3Grid()
	want := map[string]string{
		"B1": "Wyvern", "B2": "Ghost", "C4": "Troll",
		"C6": "Troll", "E2": "Ghost", "E7": "Troll", "F8": "Troll",
	}
	for code, monster := range want {
		room := w.Rooms[level3Room(code)]
		if room == nil {
			t.Fatalf("Level3Grid is missing cell %s", code)
		}
		if room.Monster != monster {
			t.Errorf("room %s Monster = %q, want %q", code, room.Monster, monster)
		}
		if room.MonsterHealth <= 0 {
			t.Errorf("room %s MonsterHealth = %d, want > 0", code, room.MonsterHealth)
		}
	}
}

// TestLevel3GridHasMantisCharm pins the first Charm-item placement in
// the whole port: Belezbar's confirmed Charm, "Mantis", cross-referenced
// from the numbered map poster into the Gorburg zone's first cell (see
// Level3Grid's doc comment for the zone-vs-exact-cell honesty caveat).
func TestLevel3GridHasMantisCharm(t *testing.T) {
	w := Level3Grid()
	room := w.Rooms[level3Room("A1")]
	found := false
	for _, item := range room.Items {
		if item == "Mantis" {
			found = true
		}
	}
	if !found {
		t.Errorf("Level3Grid room A1 Items = %v, want it to include \"Mantis\"", room.Items)
	}
}

// TestLevel3GridA2HasPelletSwap pins round 133's placement: A2, in the
// same GORBURG zone already used for Belezbar's Mantis (A1), carries
// the real "protected item" swap mechanic (see world.Room.SwapItem's
// doc comment) - a Ball dropped here reveals a Pellet.
func TestLevel3GridA2HasPelletSwap(t *testing.T) {
	w := Level3Grid()
	room := w.Rooms[level3Room("A2")]
	if room == nil || room.SwapItem != "Ball" || room.RevealItem != "Pellet" {
		t.Fatalf("Level3Grid room A2 SwapItem/RevealItem = %q/%q, want \"Ball\"/\"Pellet\"", room.SwapItem, room.RevealItem)
	}
}

// TestLevel3GridA3IsGorburgWithChest pins round 179's addition: a full
// frame-by-frame review of real gameplay footage shows the real, live
// status-panel name "Gorburg" at a room with a real, examinable chest
// holding a Leaf and a Bag - the same zone already anchoring A1's
// Mantis and A2's Pellet swap.
func TestLevel3GridA3IsGorburgWithChest(t *testing.T) {
	w := Level3Grid()
	room := w.Rooms[level3Room("A3")]
	if room == nil || room.Name != "Gorburg" || !room.HasChest {
		t.Fatalf("Level3Grid room A3 = %+v, want Name \"Gorburg\" with a real chest", room)
	}
	for _, want := range []string{"Leaf", "Bag"} {
		found := false
		for _, item := range room.Items {
			if item == want {
				found = true
			}
		}
		if !found {
			t.Errorf("Gorburg (A3) Items = %v, want it to include %q", room.Items, want)
		}
	}
}

// TestLevel3GridMainComponentIsFullyConnected pins that the original
// 41-cell component, like Level2Grid's 50 (but unlike Level1Grid's
// 44/64), is entirely reachable from the start room - it was
// deliberately built from one whole connected component, not assembled
// from a partially-connected automated pass. The 42nd cell (Sothic
// Complex, D4) is a real, later-added, deliberately isolated named room
// (see TestLevel3GridSothicComplexIsIsolated) and is correctly NOT part
// of this count.
func TestLevel3GridMainComponentIsFullyConnected(t *testing.T) {
	w := Level3Grid()
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
	if len(visited) != 41 {
		t.Errorf("only %d of 41 cells reachable from the start room, want all 41", len(visited))
	}
}
