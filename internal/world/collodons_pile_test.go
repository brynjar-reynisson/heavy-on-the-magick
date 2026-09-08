package world

import "testing"

// TestCollodonsPileHasFifteenRooms: 13 real rooms from the CASA
// walkthrough's own path, plus Furnace Room - a real room too (see
// CollodonsPile's doc comment), but sourced differently (a first-hand
// playthrough account, not the walkthrough's own path) and reached only
// via a failed INVOKE's real punishment teleport, not a normal
// directional exit - plus a real "Exit" room (Pile Collodom's own real
// North exit, confirmed via a full walkthrough video - see its own doc
// comment), CollodonsPile's first real, walkable win condition. Room of
// Misery's real West exit (see its own doc comment) leads to the
// EXISTING Sothic Complex room, not an extra room of its own - a real
// frame-by-frame video review corrected an earlier guess that it was a
// separate "Sign" room.
func TestCollodonsPileHasFifteenRooms(t *testing.T) {
	w := CollodonsPile()
	if len(w.Rooms) != 15 {
		t.Errorf("len(w.Rooms) = %d, want 15 (13 from the walkthrough path + Furnace Room + Exit)", len(w.Rooms))
	}
}

// TestCollodonsPileRoomOfMiseryHasRealWestExit covers a real exit
// confirmed via a direct frame-by-frame review of real gameplay footage
// (a Let's Play video showing "EXITS: W E" in Room of Misery, then "YOU
// ARE IN THE SOTHIC COMPLEX" after walking West) - see Room of Misery's
// own doc comment for the correction history (an earlier guess, from
// Level2Grid's own "Sign"-labeled F3 cell, wrongly modeled this as a
// separate room). Round-trip: the same footage shows walking back East
// returns to Room of Misery.
func TestCollodonsPileRoomOfMiseryHasRealWestExit(t *testing.T) {
	w := CollodonsPile()
	dest, ok := w.Rooms[roomMisery].Exits[West]
	if !ok {
		t.Fatal("Room of Misery has no West exit, want one leading to Sothic Complex")
	}
	if dest != roomSothicComplex {
		t.Errorf("Room of Misery's West exit leads to room ID %v, want roomSothicComplex", dest)
	}
	back, ok := w.Rooms[roomSothicComplex].Exits[East]
	if !ok || back != roomMisery {
		t.Errorf("Sothic Complex's East exit = (%v, %v), want a real return to Room of Misery", back, ok)
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

// TestCollodonsPileMethosHasErlstone pins the round-107 addition
// (Asmodee's confirmed Charm) - see CollodonsPile's doc comment for the
// numbered-map-plus-zone-banner sourcing.
func TestCollodonsPileMethosHasErlstone(t *testing.T) {
	w := CollodonsPile()
	room := w.Rooms[roomMethos]
	found := false
	for _, item := range room.Items {
		if item == "Erlstone" {
			found = true
		}
	}
	if !found {
		t.Errorf("Methos Items = %v, want it to include \"Erlstone\"", room.Items)
	}
}

// TestCollodonsPileTrollwyndHasMirror pins the round-177 addition: a
// full frame-by-frame review of real gameplay footage shows the player
// picking up a real "Mirror" item at Trollwynd, alongside the
// already-placed Clasp/Nougat/Scroll.
func TestCollodonsPileTrollwyndHasMirror(t *testing.T) {
	w := CollodonsPile()
	room := w.Rooms[roomTrollwynd]
	found := false
	for _, item := range room.Items {
		if item == "Mirror" {
			found = true
		}
	}
	if !found {
		t.Errorf("Trollwynd Items = %v, want it to include \"Mirror\"", room.Items)
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

// TestCollodonsPileWolfdorpHasWerewolf pins the round-120 addition:
// zone_monsters.go's "Wolfdorp: Ghost x2, Werewolf x2" sighting had
// sat unapplied - Werewolf was chosen (not Ghost) because it activates
// the real, already-tested game.checkNougatWerewolf mechanic for the
// first time in actual default-mode play. See CollodonsPile's doc
// comment.
func TestCollodonsPileWolfdorpHasWerewolf(t *testing.T) {
	w := CollodonsPile()
	room := w.Rooms[roomWolfdorp]
	if room.Monster != "Werewolf" || room.MonsterHealth <= 0 {
		t.Errorf("Wolfdorp Monster = %q (health %d), want a live Werewolf", room.Monster, room.MonsterHealth)
	}
}

// TestCollodonsPileSecundaPortaHasSign pins round 122: the manual's own
// opening narrative ("In the next room was a Sign . . .") describes the
// room right after Room of Misery, which CollodonsPile's own confirmed
// walkthrough path already establishes as Secunda Porta. See
// CollodonsPile's doc comment.
func TestCollodonsPileSecundaPortaHasSign(t *testing.T) {
	w := CollodonsPile()
	room := w.Rooms[roomSecundaPorta]
	found := false
	for _, item := range room.Items {
		if item == "Sign" {
			found = true
		}
	}
	if !found {
		t.Errorf("Secunda Porta Items = %v, want it to include \"Sign\"", room.Items)
	}
}

// TestCollodonsPileMethosHasWraith pins the round-71 addition: Methos
// is a real, connected, reachable room, and zone_monsters.go's
// independently-sourced "Methos: Wraith x1" sighting - cross-validated
// against other exact matches in that same list - is the first monster
// this room has had (see CollodonsPile's doc comment). Renamed from
// "Wraith" to "Vampire" in round 74 (see Level1Grid's doc comment), then
// reverted back to "Wraith" in round 177: a full frame-by-frame review
// of real gameplay footage shows Methos's own live combat text reads
// "WRAITH IS DEAD", while a separate encounter (Morfang) shows real
// combat text "VAMPIRE ATTACKS!" for what's now confirmed to be a
// genuinely different monster - see CollodonsPile's doc comment.
func TestCollodonsPileMethosHasWraith(t *testing.T) {
	w := CollodonsPile()
	room := w.Rooms[roomMethos]
	if room.Monster != "Wraith" || room.MonsterHealth <= 0 {
		t.Errorf("Methos Monster = %q (health %d), want a live Wraith", room.Monster, room.MonsterHealth)
	}
}

// TestCollodonsPileMethosHasCauldronBones pins the round-177 addition:
// a full frame-by-frame review of real gameplay footage shows the
// player picking up Ulna, Thigh, and Skull at Methos - the exact 3
// ingredients game.cauldronAchad's "CAULDRON, ACHAD" ritual (round 139)
// requires, previously sourced but never placed anywhere reachable.
func TestCollodonsPileMethosHasCauldronBones(t *testing.T) {
	w := CollodonsPile()
	room := w.Rooms[roomMethos]
	for _, want := range []string{"Ulna", "Thigh", "Skull"} {
		found := false
		for _, item := range room.Items {
			if item == want {
				found = true
			}
		}
		if !found {
			t.Errorf("Methos Items = %v, want it to include %q", room.Items, want)
		}
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

// TestCollodonsPileDoorHintsMatchRealPasswords covers round 159's real,
// sourced riddle text (see Room.DoorHints's doc comment) - Wolfdorp and
// Pilefoot both have real DoorHints, and every room's hint count should
// never exceed its own real password count (a hint with no matching
// password would be a fabrication this project's own discipline
// wouldn't allow).
func TestCollodonsPileDoorHintsMatchRealPasswords(t *testing.T) {
	w := CollodonsPile()
	for _, id := range []RoomID{roomWolfdorp, roomPilefoot} {
		room := w.Rooms[id]
		if len(room.DoorHints) == 0 {
			t.Errorf("room %q DoorHints is empty, want real riddle text", room.Name)
		}
		if len(room.DoorHints) > len(room.DoorPasswords) {
			t.Errorf("room %q has %d DoorHints but only %d DoorPasswords - a hint with no matching password", room.Name, len(room.DoorHints), len(room.DoorPasswords))
		}
	}
}
