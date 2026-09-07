package world

import "testing"

func TestCollodonsPileHasThirteenRooms(t *testing.T) {
	w := CollodonsPile()
	if len(w.Rooms) != 13 {
		t.Errorf("len(w.Rooms) = %d, want 13 (per the walkthrough this data was sourced from)", len(w.Rooms))
	}
}

func TestCollodonsPileStartsInRoomOfMisery(t *testing.T) {
	w := CollodonsPile()
	room := w.CurrentRoom()
	if room == nil || room.Name != "Room of Misery" {
		t.Fatalf("starting room = %+v, want Room of Misery", room)
	}
}

func TestCollodonsPileWalkthroughPath(t *testing.T) {
	// Confirms the full sourced path from the walkthrough can actually be
	// walked end to end: Room of Misery -E-> Secunda Porta -N-> Trollwynd
	// -N-> Agile Stair -SE-> Methos -S-> Sothic Complex -S-> Wolfdorp
	// -NW-> Room of Stings -N-> Morfang -E-> Room of Arrows -E-> Nidus
	// -W-> Pilefoot -N-> Pile Collodom.
	w := CollodonsPile()
	path := []struct {
		dir      Direction
		wantName string
	}{
		{East, "Secunda Porta"},
		{North, "Trollwynd"},
		{North, "Agile Stair"},
		{SouthEast, "Methos"},
		{South, "Sothic Complex"},
		{South, "Wolfdorp"},
		{NorthWest, "Room of Stings"},
		{North, "Morfang"},
		{East, "Room of Arrows"},
		{East, "Nidus"},
		{West, "Pilefoot"},
		{North, "Pile Collodom"},
	}
	for _, step := range path {
		if !w.Move(step.dir) {
			t.Fatalf("Move(%v) failed, expected to reach %q", step.dir, step.wantName)
		}
		got := w.CurrentRoom()
		if got == nil || got.Name != step.wantName {
			t.Fatalf("after Move(%v), current room = %+v, want %q", step.dir, got, step.wantName)
		}
	}

	if len(w.VisitedRooms()) != 13 {
		t.Errorf("after walking the full path, len(VisitedRooms()) = %d, want 13", len(w.VisitedRooms()))
	}
}

func TestCollodonsPileSecondaryExits(t *testing.T) {
	// The two additional, directly-stated (not inferred) exits found in a
	// second extraction pass over the same walkthrough.
	w := CollodonsPile()

	w.Current = roomTrollwynd
	if !w.Move(South) {
		t.Fatal("expected Trollwynd to have a South exit to Sothic Complex")
	}
	if got := w.CurrentRoom().Name; got != "Sothic Complex" {
		t.Errorf("Trollwynd -South-> %q, want Sothic Complex", got)
	}

	w.Current = roomArrows
	if !w.Move(North) {
		t.Fatal("expected Room of Arrows to have a North exit to Wolfdorp")
	}
	if got := w.CurrentRoom().Name; got != "Wolfdorp" {
		t.Errorf("Room of Arrows -North-> %q, want Wolfdorp", got)
	}
}

// TestCollodonsPileWolfdorpHasSword pins the real, cross-referenced
// Sword placement (Astarot's confirmed Charm) - see CollodonsPile's doc
// comment for the two-source sourcing.
func TestCollodonsPileWolfdorpHasSword(t *testing.T) {
	w := CollodonsPile()
	room := w.Rooms[roomWolfdorp]
	found := false
	for _, item := range room.Items {
		if item == "Sword" {
			found = true
		}
	}
	if !found {
		t.Errorf("Wolfdorp Items = %v, want it to include \"Sword\"", room.Items)
	}
}

// TestCollodonsPileWolfdorpDoesNotHaveKey pins the round-82 correction:
// round 81 placed "Key" in Wolfdorp's Items based on a misread
// AI-summarized list, not the raw source text - the raw text shows the
// real pickup happens in an unnamed room reached by leaving Wolfdorp
// (the same "unnamed intermediate room" trap round 64 already caught
// once for Pilefoot). See CollodonsPile's doc comment for the full
// correction writeup. Room of Stings' TollItem "Key" is honestly
// unplaced/unsourced within CollodonsPile again, as it was from round
// 64 through round 80.
func TestCollodonsPileWolfdorpDoesNotHaveKey(t *testing.T) {
	w := CollodonsPile()
	room := w.Rooms[roomWolfdorp]
	for _, item := range room.Items {
		if item == "Key" {
			t.Errorf("Wolfdorp Items = %v, want it to NOT include \"Key\" (round-82 correction - see doc comment)", room.Items)
		}
	}
}

// TestCollodonsPileSothicComplexHasSunflower pins the real, cross-
// referenced Sunflower placement (Magot's confirmed Charm) - see
// CollodonsPile's doc comment for the numbered-map-plus-zone-banner
// sourcing (same method as Wolfdorp's Sword).
func TestCollodonsPileSothicComplexHasSunflower(t *testing.T) {
	w := CollodonsPile()
	room := w.Rooms[roomSothicComplex]
	found := false
	for _, item := range room.Items {
		if item == "Sunflower" {
			found = true
		}
	}
	if !found {
		t.Errorf("Sothic Complex Items = %v, want it to include \"Sunflower\"", room.Items)
	}
}

// TestCollodonsPileRoomOfMiseryHasBothNumberedItems pins Room of
// Misery's 2 items - the numbered map poster labels this exact room
// "1, 2", the strongest-confidence item placement in this file (see
// CollodonsPile's doc comment).
func TestCollodonsPileRoomOfMiseryHasBothNumberedItems(t *testing.T) {
	w := CollodonsPile()
	room := w.Rooms[roomMisery]
	want := map[string]bool{"Grimoire": false, "Poison-smeared book": false}
	for _, item := range room.Items {
		if _, ok := want[item]; ok {
			want[item] = true
		}
	}
	for item, found := range want {
		if !found {
			t.Errorf("Room of Misery Items = %v, want it to include %q", room.Items, item)
		}
	}
}

// TestCollodonsPileHasTablePlacements pins the 8 real rooms where the
// CASA walkthrough confirms "EXAMINE TABLE" (see world.Room.HasTable's
// doc comment).
func TestCollodonsPileHasTablePlacements(t *testing.T) {
	w := CollodonsPile()
	want := []RoomID{roomMisery, roomTrollwynd, roomSothicComplex, roomWolfdorp, roomStings, roomMorfang, roomArrows, roomMethos}
	for _, id := range want {
		room := w.Rooms[id]
		if room == nil || !room.HasTable {
			t.Errorf("room %q HasTable = %v, want true", room.Name, room != nil && room.HasTable)
		}
	}
}

// TestCollodonsPileHasChestPlacements pins the 2 real rooms where the
// CASA walkthrough confirms "EXAMINE CHEST" (see world.Room.HasChest's
// doc comment) - a distinct container from HasTable's "EXAMINE TABLE".
func TestCollodonsPileHasChestPlacements(t *testing.T) {
	w := CollodonsPile()
	want := []RoomID{roomWolfdorp, roomMorfang}
	for _, id := range want {
		room := w.Rooms[id]
		if room == nil || !room.HasChest {
			t.Errorf("room %q HasChest = %v, want true", room.Name, room != nil && room.HasChest)
		}
	}
}

// TestCollodonsPileMethosHasVampire pins the round-71 addition: Methos
// is a real, connected, reachable room, and zone_monsters.go's
// independently-sourced "Methos: Wraith x1" sighting - cross-validated
// against other exact matches in that same list - is the first monster
// this room has had (see CollodonsPile's doc comment). Renamed from
// "Wraith" to "Vampire" in round 74 (see Level1Grid's doc comment).
func TestCollodonsPileMethosHasVampire(t *testing.T) {
	w := CollodonsPile()
	room := w.Rooms[roomMethos]
	if room.Monster != "Vampire" || room.MonsterHealth <= 0 {
		t.Errorf("Methos Monster = %q (health %d), want a live Vampire", room.Monster, room.MonsterHealth)
	}
}

// TestCollodonsPileMorfangHasVampire pins the round-106 addition:
// zone_monsters.go's "Morfang: Vampire x3" sighting was left unapplied
// for many rounds because its COUNT doesn't match Level1Grid's 4
// shipped Vampires (F2/G1/G2/H1) - but CollodonsPile's single Monster
// field only records species, not count, and both sources agree
// unambiguously on that. See CollodonsPile's doc comment.
func TestCollodonsPileMorfangHasVampire(t *testing.T) {
	w := CollodonsPile()
	room := w.Rooms[roomMorfang]
	if room.Monster != "Vampire" || room.MonsterHealth <= 0 {
		t.Errorf("Morfang Monster = %q (health %d), want a live Vampire", room.Monster, room.MonsterHealth)
	}
}

// TestCollodonsPileTrollwyndAndSothicComplexHaveScrollNougat pins the
// round-63 walkthrough re-read: Nougat and a Scroll in Trollwynd, plus
// a second, separate Scroll in Sothic Complex - resolving 2 of the 3
// items an earlier pass had left unplaced as "vague/multiple locations"
// (see CollodonsPile's doc comment).
func TestCollodonsPileTrollwyndAndSothicComplexHaveScrollNougat(t *testing.T) {
	w := CollodonsPile()
	want := map[RoomID][]string{
		roomTrollwynd:     {"Nougat", "Scroll"},
		roomSothicComplex: {"Scroll"},
	}
	for roomID, items := range want {
		room := w.Rooms[roomID]
		for _, want := range items {
			found := false
			for _, item := range room.Items {
				if item == want {
					found = true
				}
			}
			if !found {
				t.Errorf("room %q Items = %v, want it to include %q", room.Name, room.Items, want)
			}
		}
	}
}
