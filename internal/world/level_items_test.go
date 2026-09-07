package world

import "testing"

func TestLevelItemsNonEmptyAndTagged(t *testing.T) {
	for _, items := range [][]LevelItem{LevelOneItems, LevelTwoItems, LevelThreeItems, LevelFourItems} {
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

// TestLevelFourItemsCrossConfirmRabak pins round 149's find: this
// poster's own "heavymap-levels3-4-poster.jpg" half independently
// names "Rabak" on Level Four, matching level4_grid.go's own already-
// placed "Doubt of Rabak" special room (D3) - a real cross-source
// confirmation, the same pattern TestLevelItemsCrossConfirmPileCollodom
// already established for Level One/Pile Collodom.
func TestLevelFourItemsCrossConfirmRabak(t *testing.T) {
	for _, it := range LevelFourItems {
		if it.Name == "Rabak" {
			return
		}
	}
	t.Error(`LevelFourItems should include "Rabak", cross-confirming "Doubt of Rabak" in level4_grid.go`)
}
