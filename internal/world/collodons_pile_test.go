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
