package game

import "fmt"

// combatStaminaCost is charged to the player per BLAST/FREEZE cast against
// a real monster - a placeholder amount, not extracted (see Handle's doc
// comment for the sourcing of the underlying "combat costs Stamina" fact).
const combatStaminaCost = 5

// transfusion handles TRANSFUSION. Confirmed real effect (restores
// Stamina) via the CASA walkthrough extraction; the exact original
// restore amount was not stated there and hasn't been extracted from the
// disassembly, so the amount itself is an honest placeholder. The cap at
// character.Player.MaxStamina IS a real, sourced mechanic though: the
// same walkthrough instructs casting TRANSFUSION "a few times until you
// get maximum stamina" - wording that only makes sense if repeated casts
// approach a real ceiling rather than healing indefinitely.
func (g *Game) transfusion() string {
	const placeholderRestoreAmount = 10
	g.Player.Stamina = min(g.Player.Stamina+placeholderRestoreAmount, g.Player.MaxStamina)
	if g.Player.Stamina >= g.Player.MaxStamina {
		return fmt.Sprintf("You feel your strength return. Stamina is now %d (maximum). (restore-per-cast amount is a placeholder - not yet extracted from the original)", g.Player.Stamina)
	}
	return fmt.Sprintf("You feel your strength return. Stamina is now %d. (restore-per-cast amount is a placeholder - not yet extracted from the original)", g.Player.Stamina)
}

func (g *Game) blast() string {
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
func (g *Game) deathCheck(msg string) string {
	if g.Player.IsDead() {
		return msg + "\nYour Stamina gives out. You are dead. (GAME OVER)"
	}
	return msg
}
