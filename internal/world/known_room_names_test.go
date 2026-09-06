package world

import "testing"

func TestKnownUnplacedRoomNamesNonEmpty(t *testing.T) {
	if len(KnownUnplacedRoomNames) == 0 {
		t.Fatal("expected at least one known-but-unplaced room name")
	}
	for _, name := range KnownUnplacedRoomNames {
		if name == "" {
			t.Error("KnownUnplacedRoomNames contains an empty name")
		}
	}
}
