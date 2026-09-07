package magic

import "testing"

func TestDemonsConfirmedFour(t *testing.T) {
	if len(Demons) != 4 {
		t.Fatalf("len(Demons) = %d, want 4 (per the official manual's grimoire section)", len(Demons))
	}
	names := map[string]bool{}
	for _, d := range Demons {
		names[d.Name] = true
		if d.Title == "" || d.Number == 0 || d.Charm == "" || d.Correspondences == "" {
			t.Errorf("Demon %+v missing expected fields", d)
		}
	}
	for _, want := range []string{"ASMODEE", "ASTAROT", "BELEZBAR", "MAGOT"} {
		if !names[want] {
			t.Errorf("Demons missing %q", want)
		}
	}
}

func TestDemonAbilitiesKnownForTwo(t *testing.T) {
	ability := map[string]string{}
	for _, d := range Demons {
		ability[d.Name] = d.Ability
	}
	if ability["ASTAROT"] == "" {
		t.Error("Astarot should have a confirmed Ability (transportation, per the map's demon panel)")
	}
	if ability["MAGOT"] == "" {
		t.Error("Magot should have a confirmed Ability (locates objects, per the map's demon panel)")
	}
}
