package world

import "testing"

func TestNumberedRoomContentsComplete(t *testing.T) {
	if len(NumberedRoomContents) != 102 {
		t.Fatalf("len(NumberedRoomContents) = %d, want 102", len(NumberedRoomContents))
	}
	for i, entry := range NumberedRoomContents {
		wantNumber := i + 1
		if entry.Number != wantNumber {
			t.Errorf("entry %d has Number %d, want %d (list should be in order, no gaps)", i, entry.Number, wantNumber)
		}
		if entry.Contents == "" {
			t.Errorf("entry #%d has empty Contents", entry.Number)
		}
	}
}
