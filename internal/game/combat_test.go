package game

import (
	"strings"
	"testing"

	"github.com/brynjar-reynisson/heavy-on-the-magick/internal/parser"
)

// withGrimoire grants g's player the real Grimoire item directly (rather
// than walking to Room of Misery and picking it up), simulating "already
// found it" for tests that exercise BLAST/FREEZE/TRANSFUSION - see
// spellRequiresItem's doc comment for why this is now required.
func withGrimoire(g *Game) {
	g.Player.Items = append(g.Player.Items, "Grimoire")
}

func TestHandleAttackAndKillAreSynonymsForBlast(t *testing.T) {
	for _, verb := range []string{"ATTACK", "ATTACKS", "KILL"} {
		g := New() // Room of Misery has no known Monster
		withGrimoire(g)
		got := g.Handle(parser.Parse(verb))
		if !strings.Contains(got, "nothing here") {
			t.Errorf("Handle(%s) with no monster = %q, want the same response as BLAST", verb, got)
		}
	}
}

func TestHandleBlastWithNoMonster(t *testing.T) {
	g := New() // Room of Misery has no known Monster
	withGrimoire(g)
	got := g.Handle(parser.Parse("BLAST"))
	if !strings.Contains(got, "nothing here") {
		t.Errorf("Handle(BLAST) with no monster = %q, want it to say there's nothing to hit", got)
	}
}

// TestHandleBlastRequiresGrimoire covers the real gate found via a Let's
// Play video of the original (see spellRequiresItem's doc comment): the
// game's own real rejection ("YOU CAN'T INVOKE SPELL") appears before
// Axil has picked up the Grimoire, even with a real monster present.
func TestHandleBlastRequiresGrimoire(t *testing.T) {
	g := New()
	g.Handle(parser.Parse("EAST"))
	g.Handle(parser.Parse("DOOR, SILENCE"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd, has a Monster
	room := g.World.CurrentRoom()
	startHealth := room.MonsterHealth

	got := g.Handle(parser.Parse("BLAST"))
	if !strings.Contains(got, "Grimoire") {
		t.Errorf("Handle(BLAST) without the Grimoire = %q, want a rejection naming the Grimoire", got)
	}
	if room.MonsterHealth != startHealth {
		t.Errorf("MonsterHealth changed to %d despite the missing Grimoire, want unchanged %d", room.MonsterHealth, startHealth)
	}
}

// TestHandleFreezeRequiresGrimoire and
// TestHandleTransfusionRequiresGrimoire cover the same real gate for the
// other 2 spells it applies to (see spellRequiresItem).
func TestHandleFreezeRequiresGrimoire(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("FREEZE"))
	if !strings.Contains(got, "Grimoire") {
		t.Errorf("Handle(FREEZE) without the Grimoire = %q, want a rejection naming the Grimoire", got)
	}
}

func TestHandleTransfusionRequiresGrimoire(t *testing.T) {
	g := New()
	got := g.Handle(parser.Parse("TRANSFUSION"))
	if !strings.Contains(got, "Grimoire") {
		t.Errorf("Handle(TRANSFUSION) without the Grimoire = %q, want a rejection naming the Grimoire", got)
	}
}

func TestHandleBlastDefeatsMonster(t *testing.T) {
	g := New()
	withGrimoire(g)
	g.Handle(parser.Parse("EAST")) // Secunda Porta
	g.Handle(parser.Parse("DOOR, SILENCE"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd, which has a Monster

	room := g.World.CurrentRoom()
	if room.Name != "Trollwynd" || room.Monster == "" {
		t.Fatalf("test setup bug: expected to be in Trollwynd with a monster, got %+v", room)
	}
	startHealth := room.MonsterHealth

	// BLAST's per-hit damage now scales with Skill (see blastDamage), so
	// the monster may take fewer hits than its starting MonsterHealth to
	// defeat - loop until it's actually dead rather than assuming 1
	// damage per hit, with a generous safety cap against an infinite loop.
	var got string
	hits := 0
	for room.MonsterHealth > 0 && hits < startHealth+5 {
		got = g.Handle(parser.Parse("BLAST"))
		hits++
		if room.MonsterHealth > 0 && !strings.Contains(got, "still standing") {
			t.Errorf("Handle(BLAST) mid-fight = %q, want it to say the monster is still standing", got)
		}
	}
	if room.MonsterHealth > 0 {
		t.Fatalf("MonsterHealth after %d BLASTs = %d, want <= 0 (defeated)", hits, room.MonsterHealth)
	}
	if !strings.Contains(got, "destroyed") {
		t.Errorf("Handle(BLAST) on the killing blow = %q, want it to say destroyed", got)
	}

	// One more BLAST after defeat should report nothing left to hit.
	got = g.Handle(parser.Parse("BLAST"))
	if !strings.Contains(got, "nothing here") {
		t.Errorf("Handle(BLAST) after monster defeated = %q, want nothing-to-hit response", got)
	}
}

func TestHandleFreezeAwardsExperiencePoints(t *testing.T) {
	g := New()
	withGrimoire(g)
	g.Handle(parser.Parse("EAST"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd
	before := g.Player.ExperiencePoints

	got := g.Handle(parser.Parse("FREEZE"))
	want := 10 + g.Player.Luck
	if g.Player.ExperiencePoints != before+want {
		t.Errorf("ExperiencePoints after defeating a monster = %d, want %d (before %d + base 10 + Luck %d)", g.Player.ExperiencePoints, before+want, before, g.Player.Luck)
	}
	if !strings.Contains(got, "Experience Points") {
		t.Errorf("Handle(FREEZE) on a kill = %q, want it to mention Experience Points", got)
	}
}

func TestHandleBlastDefeatsMonsterAwardsExperiencePoints(t *testing.T) {
	g := New()
	withGrimoire(g)
	g.Handle(parser.Parse("EAST"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd
	room := g.World.CurrentRoom()
	before := g.Player.ExperiencePoints

	var got string
	hits := 0
	for room.MonsterHealth > 0 && hits < 10 {
		got = g.Handle(parser.Parse("BLAST"))
		hits++
	}
	if g.Player.ExperiencePoints <= before {
		t.Errorf("ExperiencePoints after defeating a monster with BLAST = %d, want more than before (%d)", g.Player.ExperiencePoints, before)
	}
	if !strings.Contains(got, "Experience Points") {
		t.Errorf("Handle(BLAST) on the killing blow = %q, want it to mention Experience Points", got)
	}
}

func TestHandleFreezeDefeatsMonsterInstantly(t *testing.T) {
	g := New()
	withGrimoire(g)
	g.Handle(parser.Parse("EAST"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd

	got := g.Handle(parser.Parse("FREEZE"))
	if !strings.Contains(got, "FREEZE") {
		t.Errorf("Handle(FREEZE) = %q, want it to mention freezing", got)
	}
	room := g.World.CurrentRoom()
	if room.MonsterHealth != 0 {
		t.Errorf("MonsterHealth after FREEZE = %d, want 0 (instant neutralize)", room.MonsterHealth)
	}
}

func TestBlastDamageScalesWithSkill(t *testing.T) {
	g := New()
	g.Player.Skill = 4 // bottom of the roll range
	if got, want := g.blastDamage(), 2; got != want {
		t.Errorf("blastDamage() with Skill 4 = %d, want %d", got, want)
	}
	g.Player.Skill = 12 // top of the roll range
	if got, want := g.blastDamage(), 4; got != want {
		t.Errorf("blastDamage() with Skill 12 = %d, want %d", got, want)
	}
}

func TestHandleBlastCostsStamina(t *testing.T) {
	g := New()
	withGrimoire(g)
	g.Handle(parser.Parse("EAST"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd, has a Monster
	before := g.Player.Stamina
	g.Handle(parser.Parse("BLAST"))
	if g.Player.Stamina != before-combatStaminaCost {
		t.Errorf("Stamina after BLAST = %d, want %d (before %d minus combatStaminaCost %d)", g.Player.Stamina, before-combatStaminaCost, before, combatStaminaCost)
	}
}

func TestHandleFreezeCostsStamina(t *testing.T) {
	g := New()
	withGrimoire(g)
	g.Handle(parser.Parse("EAST"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd, has a Monster
	before := g.Player.Stamina
	g.Handle(parser.Parse("FREEZE"))
	if g.Player.Stamina != before-combatStaminaCost {
		t.Errorf("Stamina after FREEZE = %d, want %d (before %d minus combatStaminaCost %d)", g.Player.Stamina, before-combatStaminaCost, before, combatStaminaCost)
	}
}

func TestHandleCombatCanKillPlayer(t *testing.T) {
	g := New()
	withGrimoire(g)
	g.Handle(parser.Parse("EAST"))
	g.Handle(parser.Parse("NORTH")) // Trollwynd, has a Monster
	g.Player.Stamina = combatStaminaCost

	got := g.Handle(parser.Parse("BLAST"))
	if !g.Player.IsDead() {
		t.Fatalf("Player.Stamina = %d after BLAST, want <= 0 (dead)", g.Player.Stamina)
	}
	if !strings.Contains(got, "die horribly") {
		t.Errorf("Handle(BLAST) that reduces Stamina to 0 = %q, want the real confirmed death message", got)
	}

	got = g.Handle(parser.Parse("EAST"))
	if got != "You are dead. (Stamina reached 0 - GAME OVER)" {
		t.Errorf("Handle(...) after death = %q, want the dead-rejection message", got)
	}
}

func TestHandleDeadPlayerCanStillLookAndMap(t *testing.T) {
	g := New()
	g.Player.Stamina = 0
	if got := g.Handle(parser.Parse("LOOK")); strings.Contains(got, "You are dead") {
		t.Errorf("Handle(LOOK) while dead = %q, want LOOK to still work", got)
	}
	if got := g.Handle(parser.Parse("MAP")); strings.Contains(got, "You are dead") {
		t.Errorf("Handle(MAP) while dead = %q, want MAP to still work", got)
	}
}

func TestHandleTransfusionRestoresStamina(t *testing.T) {
	g := New()
	withGrimoire(g)
	g.Player.Stamina -= 5 // take damage first; a fresh player starts at MaxStamina already
	before := g.Player.Stamina
	got := g.Handle(parser.Parse("TRANSFUSION"))
	if g.Player.Stamina <= before {
		t.Errorf("Stamina after TRANSFUSION = %d, want more than before (%d)", g.Player.Stamina, before)
	}
	if !strings.Contains(got, "Stamina") {
		t.Errorf("Handle(TRANSFUSION) = %q, want it to report Stamina", got)
	}
}

func TestHandleTransfusionCapsAtMaxStamina(t *testing.T) {
	g := New() // a fresh player already starts at MaxStamina
	withGrimoire(g)
	got := g.Handle(parser.Parse("TRANSFUSION"))
	if g.Player.Stamina != g.Player.MaxStamina {
		t.Errorf("Stamina after TRANSFUSION on a full-health player = %d, want it capped at MaxStamina %d", g.Player.Stamina, g.Player.MaxStamina)
	}
	if !strings.Contains(got, "maximum") {
		t.Errorf("Handle(TRANSFUSION) at full health = %q, want it to say maximum", got)
	}
}
