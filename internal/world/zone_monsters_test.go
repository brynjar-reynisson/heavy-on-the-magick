package world

import "testing"

func TestZoneMonsterSightingsWellFormed(t *testing.T) {
	if len(ZoneMonsterSightings) == 0 {
		t.Fatal("expected at least one zone monster sighting")
	}
	for _, s := range ZoneMonsterSightings {
		if s.Zone == "" || s.Monster == "" || s.Count < 1 {
			t.Errorf("ZoneMonsterSighting %+v is malformed", s)
		}
	}
}

func TestZoneMonsterSightingsCrossConfirmsNidusCyclops(t *testing.T) {
	for _, s := range ZoneMonsterSightings {
		if s.Zone == "Nidus" && s.Monster == "Cyclops" {
			return
		}
	}
	t.Error(`expected a "Nidus"/"Cyclops" sighting, cross-confirming CollodonsPile's Nidus room`)
}
