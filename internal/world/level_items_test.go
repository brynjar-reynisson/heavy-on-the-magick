package world

import "testing"

func TestLevelItemsNonEmptyAndTagged(t *testing.T) {
	for _, items := range [][]LevelItem{LevelOneItems, LevelTwoItems} {
		if len(items) == 0 {
			t.Fatal("expected a non-empty item list")
		}
		wantLevel := items[0].Level
		for _, it := range items {
			if it.Level != wantLevel {
				t.Errorf("item %+v has Level %d, want %d for this slice", it, it.Level, wantLevel)
			}
			if it.Name == "" {
				t.Errorf("item %+v has an empty Name", it)
			}
		}
	}
}

func TestLevelItemsCrossConfirmPileCollodom(t *testing.T) {
	for _, it := range LevelOneItems {
		if it.Name == "Pile Collodom" {
			return
		}
	}
	t.Error(`LevelOneItems should include "Pile Collodom", cross-confirming the room of the same name in CollodonsPile`)
}
