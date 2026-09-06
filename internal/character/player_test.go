package character

import "testing"

func TestHasItem(t *testing.T) {
	p := &Player{Items: []string{"Grimoire", "Sword"}}
	if !p.HasItem("sword") {
		t.Error("HasItem should be case-insensitive and find a carried item")
	}
	if p.HasItem("Mantis") {
		t.Error("HasItem should not find an item that isn't carried")
	}
}

func TestNewPlayerStatsWithinBounds(t *testing.T) {
	// Roll several times since stats are random - a single sample
	// wouldn't catch an off-by-one in the bounds.
	for range 50 {
		p := NewPlayer()
		if p.Stamina < minStamina || p.Stamina > maxStamina {
			t.Fatalf("Stamina = %d, want in [%d, %d]", p.Stamina, minStamina, maxStamina)
		}
		if p.Skill < minSkill || p.Skill > maxSkill {
			t.Fatalf("Skill = %d, want in [%d, %d]", p.Skill, minSkill, maxSkill)
		}
		if p.Luck < minLuck || p.Luck > maxLuck {
			t.Fatalf("Luck = %d, want in [%d, %d]", p.Luck, minLuck, maxLuck)
		}
	}
}

func TestNewPlayerStartsAtNeophyte(t *testing.T) {
	p := NewPlayer()
	if p.Grade != Neophyte {
		t.Errorf("Grade = %v, want Neophyte", p.Grade)
	}
	if p.Name != "Axil" {
		t.Errorf("Name = %q, want Axil", p.Name)
	}
}

func TestRealignStaysWithinBounds(t *testing.T) {
	p := NewPlayer()
	for range 50 {
		p.Realign()
		if p.Stamina < minStamina || p.Stamina > maxStamina {
			t.Fatalf("Stamina = %d after Realign, want in [%d, %d]", p.Stamina, minStamina, maxStamina)
		}
		if p.Skill < minSkill || p.Skill > maxSkill {
			t.Fatalf("Skill = %d after Realign, want in [%d, %d]", p.Skill, minSkill, maxSkill)
		}
		if p.Luck < minLuck || p.Luck > maxLuck {
			t.Fatalf("Luck = %d after Realign, want in [%d, %d]", p.Luck, minLuck, maxLuck)
		}
	}
}

func TestNewPlayerStartsAtFullMaxStamina(t *testing.T) {
	p := NewPlayer()
	if p.MaxStamina != p.Stamina {
		t.Errorf("MaxStamina = %d, Stamina = %d, want equal for a freshly rolled player", p.MaxStamina, p.Stamina)
	}
}

func TestRealignResetsMaxStamina(t *testing.T) {
	p := NewPlayer()
	p.Stamina = 1 // simulate damage
	p.Realign()
	if p.MaxStamina != p.Stamina {
		t.Errorf("after Realign, MaxStamina = %d, Stamina = %d, want equal", p.MaxStamina, p.Stamina)
	}
}

func TestIsDead(t *testing.T) {
	p := NewPlayer()
	if p.IsDead() {
		t.Fatal("a freshly rolled player should never start dead (Stamina is always > 0 per the roll bounds)")
	}
	p.Stamina = 0
	if !p.IsDead() {
		t.Error("IsDead() with 0 Stamina = false, want true")
	}
	p.Stamina = -3
	if !p.IsDead() {
		t.Error("IsDead() with negative Stamina = false, want true")
	}
}
