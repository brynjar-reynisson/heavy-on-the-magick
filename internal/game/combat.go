package game

import "fmt"

// combatStaminaCost is charged to the player per BLAST/FREEZE cast against
// a real monster - a placeholder amount, not extracted (see Handle's doc
// comment for the sourcing of the underlying "combat costs Stamina" fact).
const combatStaminaCost = 5

// spellRequiresItem maps each real, already-implemented action-form
// spell to the specific item Axil must be carrying to cast it at all.
// Real gameplay footage of the original (a Let's Play video) shows the
// game's own actual rejection text, "YOU CAN'T INVOKE SPELL", produced
// when the player tries to cast before picking up the Grimoire (Room of
// Misery's own starting item) - a real, previously-unmodeled gate, not
// merely an inert pickup as this port had it. Deliberately a map, not a
// single hardcoded Grimoire check: Axil is confirmed to find further
// spells later in the game, and CALL already has its own separate real
// gate (the Scroll — see call()) following this exact same shape - a
// future spell's own required item is a one-line addition here, not a
// new gating mechanism to invent.
var spellRequiresItem = map[string]string{
	"BLAST":       "Grimoire",
	"FREEZE":      "Grimoire",
	"TRANSFUSION": "Grimoire",
}

// checkSpellbook reports whether the player may currently cast verb (per
// spellRequiresItem), and if not, the real, honest rejection to show -
// echoing the confirmed real "can't invoke spell" phrasing while still
// naming the missing item, matching this project's existing Talisman-gate
// convention (see invoke's own "no suitable Talisman" message).
func (g *Game) checkSpellbook(verb string) (ok bool, rejection string) {
	item, needed := spellRequiresItem[verb]
	if !needed || g.hasItem(item) {
		return true, ""
	}
	return false, fmt.Sprintf("You can't invoke that spell - you don't have the %s.", item)
}

// transfusionExperienceCost is charged per TRANSFUSION cast - a real,
// confirmed mechanic (round 147: the official map poster's own footer
// states "TRANSFUSION = STAMINA FROM EXPERIENCE"), now doubly confirmed
// by a direct frame-by-frame review of a full walkthrough video showing
// the real game's own exact rejection text, "Not enough experience.",
// when attempting to cast without sufficient ExperiencePoints. The
// EXACT cost/ratio still isn't extracted - this placeholder amount
// matches awardVictoryPoints' own base reward (10), a reasonable,
// internally-consistent guess pending the real number.
const transfusionExperienceCost = 10

// transfusion handles TRANSFUSION. Confirmed real effect (restores
// Stamina) via the CASA walkthrough extraction; the exact original
// restore amount was not stated there and hasn't been extracted from the
// disassembly, so the amount itself is an honest placeholder. The cap at
// character.Player.MaxStamina IS a real, sourced mechanic though: the
// same walkthrough instructs casting TRANSFUSION "a few times until you
// get maximum stamina" - wording that only makes sense if repeated casts
// approach a real ceiling rather than healing indefinitely.
func (g *Game) transfusion() string {
	if ok, rejection := g.checkSpellbook("TRANSFUSION"); !ok {
		return rejection
	}
	if g.Player.ExperiencePoints < transfusionExperienceCost {
		return "Not enough experience."
	}
	g.Player.ExperiencePoints -= transfusionExperienceCost
	const placeholderRestoreAmount = 10
	g.Player.Stamina = min(g.Player.Stamina+placeholderRestoreAmount, g.Player.MaxStamina)
	if g.Player.Stamina >= g.Player.MaxStamina {
		return fmt.Sprintf("You feel your strength return. Stamina is now %d (maximum). (restore-per-cast amount is a placeholder - not yet extracted from the original)", g.Player.Stamina)
	}
	return fmt.Sprintf("You feel your strength return. Stamina is now %d. (restore-per-cast amount is a placeholder - not yet extracted from the original)", g.Player.Stamina)
}

func (g *Game) blast() string {
	if ok, rejection := g.checkSpellbook("BLAST"); !ok {
		return rejection
	}
	room := g.World.CurrentRoom()
	if room == nil || room.Monster == "" || room.MonsterHealth <= 0 {
		return "You BLAST, but there's nothing here to hit."
	}
	g.Player.Stamina -= combatStaminaCost
	room.MonsterHealth -= g.blastDamage()
	if room.MonsterHealth <= 0 {
		gained := g.awardVictoryPoints()
		return g.deathCheck(fmt.Sprintf("The %s is destroyed! (+%d Experience Points)", room.Monster, gained))
	}
	return g.deathCheck(fmt.Sprintf("You BLAST the %s! It's still standing.", room.Monster))
}

// blastDamage returns how much MonsterHealth a single BLAST removes.
// Confirmed real (Spectrum Computing's plain-text instructions file for
// the game — see ../../CLAUDE.md): "your Stamina and Skill together
// affect the outcome of conflicts." No exact formula is stated, so this
// is a reasonable, documented interpretation (Skill was previously
// rolled but never consulted by any game logic at all) rather than
// extracted fact: a base 1 hit, plus 1 extra point of damage for every
// 4 points of Skill.
func (g *Game) blastDamage() int {
	return 1 + g.Player.Skill/4
}

func (g *Game) freeze() string {
	if ok, rejection := g.checkSpellbook("FREEZE"); !ok {
		return rejection
	}
	room := g.World.CurrentRoom()
	if room == nil || room.Monster == "" || room.MonsterHealth <= 0 {
		return "You FREEZE, but there's nothing here to target."
	}
	g.Player.Stamina -= combatStaminaCost
	name := room.Monster
	room.MonsterHealth = 0
	gained := g.awardVictoryPoints()
	return g.deathCheck(fmt.Sprintf("You FREEZE the %s solid! (+%d Experience Points)", name, gained))
}

// awardVictoryPoints grants ExperiencePoints for defeating a monster and
// returns the amount gained. Both the base reward and Luck's bonus are
// this port's own reasonable interpretation, not an extracted formula:
// ExperiencePoints itself is confirmed real (see character.Player's doc
// comment) but no confirmed source states what earns it or how much.
// The Luck bonus specifically has real motivation though — Spectrum
// Computing's instructions file for the game states "your Luck
// influences virtually everything you do" (see ../../CLAUDE.md), and
// Luck was previously rolled but never consulted by any game logic at
// all, same gap Skill's blastDamage bonus closed for combat damage.
func (g *Game) awardVictoryPoints() int {
	const baseVictoryPoints = 10
	gained := baseVictoryPoints + g.Player.Luck
	g.Player.ExperiencePoints += gained
	return gained
}

// deathCheck appends a death notice if the just-applied Stamina cost
// killed the player (see Handle's doc comment - confirmed real mechanic).
// "You die horribly!" is the real game's own exact death message,
// confirmed via a direct frame-by-frame review of real gameplay footage
// (a Let's Play video) - this port's own wording is kept alongside it
// for clarity (context on WHY, which the original's own short exclamation
// doesn't state) rather than replaced outright.
func (g *Game) deathCheck(msg string) string {
	if g.Player.IsDead() {
		return msg + "\nYour Stamina gives out. You die horribly! (GAME OVER)"
	}
	return msg
}
