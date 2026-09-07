package game

import (
	"fmt"
	"strings"
)

// saveStaminaCost is charged on every successful Save Game/Save Axil -
// a real, sourced mechanic (round 86): the manual states plainly "Saving
// a game will deplete your Stamina, so that a Save cannot be used as an
// easy way of getting round difficult choices!" No exact amount was
// given by the manual, so this started as an honest placeholder,
// deliberately smaller than combatStaminaCost to match the manual's own
// relative framing ("Combat will reduce your Stamina a lot, most other
// actions will reduce it a little") — round 125 found the EXACT real
// number from a direct first-hand playthrough account (The CRPG
// Addict's 2016 blog post): "You also lose 1 stamina point every time
// you save" — this already-placeholder value of 1 turns out to be
// exactly right, not a guess anymore.
const saveStaminaCost = 1

// options handles OPTIONS (Merphish keyword "O"). With no target it shows
// the Option Screen's confirmed real menu items (found live in memory
// during disassembly as the game's own on-screen UI strings under the
// "Magick!" header, not manual prose) - the disassembly confirmed these 5
// items exist among "options 1-6", but not which numbered slot each one
// occupies. Round 124 pinned down 3 of the 6 real slots from two
// separate sources: the manual states "select option 1 and Away You
// Go!" (option 1 = starting the game — no action needed here, since
// this port's game already exists once constructed) and "select option
// 6 and the values will be realigned" (option 6 = Realign Status); the
// CASA walkthrough's own "Tips" section separately states verbatim
// "SAVE regularly by pressing key O and then option 2" (option 2 =
// Save Game). Options 3-5's exact slots still aren't confirmed. A bare
// numeric target ("O 2", "O 6") now works as a real alternate way to
// trigger these, alongside the existing keyword form. "Realign Status"
// is wired to a real effect (Player.Realign, same roll ranges as
// NewPlayer) since that reroll behavior is confirmed by the manual.
// Save/Restore Game/Axil are real, confirmed menu choices (see save.go)
// — a real, functional file-based save system, though the file format
// and the Game-vs-Axil split are this port's own implementation, not the
// original's actual save mechanism (not extracted/known). Saving now
// also costs real Stamina (round 86) - the manual states plainly
// "Saving a game will deplete your Stamina, so that a Save cannot be
// used as an easy way of getting round difficult choices!" - see
// saveStaminaCost. Restoring does not cost Stamina (not stated by any
// source, and would defeat a Save's own point if it did). Round 119:
// the manual also confirms real "Version letter" save slots (see
// save.go's doc comment) - an optional trailing single-letter word in
// target (e.g. "SAVE GAME B") selects one; omitted, it's the same
// single default slot this port has always used.
func (g *Game) options(target string) string {
	target = strings.ToUpper(strings.TrimSpace(target))
	// Round 124: a bare confirmed real numeric slot (see doc comment
	// above) is translated to its equivalent keyword up front, so it
	// flows through the exact same logic below as the keyword form -
	// not a separate, duplicated code path.
	switch target {
	case "2":
		target = "SAVE GAME"
	case "6":
		target = "REALIGN"
	}
	version := extractVersionLetter(target)
	versionSuffix := ""
	if version != "" {
		versionSuffix = fmt.Sprintf(" (version %s)", version)
	}
	switch {
	case strings.Contains(target, "REALIGN"):
		g.Player.Realign()
		return fmt.Sprintf("Realign Status: Stamina %d, Skill %d, Luck %d.", g.Player.Stamina, g.Player.Skill, g.Player.Luck)
	case strings.Contains(target, "SAVE") && strings.Contains(target, "AXIL"):
		g.Player.Stamina -= saveStaminaCost
		if err := g.SaveAxilVersion(version); err != nil {
			return fmt.Sprintf("Save Axil failed: %v", err)
		}
		return g.deathCheck("Axil saved" + versionSuffix + ".")
	case strings.Contains(target, "RESTORE") && strings.Contains(target, "AXIL"):
		if err := g.RestoreAxilVersion(version); err != nil {
			return fmt.Sprintf("Restore Axil failed: %v", err)
		}
		return "Axil restored" + versionSuffix + "."
	case strings.Contains(target, "SAVE"):
		g.Player.Stamina -= saveStaminaCost
		if err := g.SaveGameVersion(version); err != nil {
			return fmt.Sprintf("Save Game failed: %v", err)
		}
		return g.deathCheck("Game saved" + versionSuffix + ".")
	case strings.Contains(target, "RESTORE"):
		if err := g.RestoreGameVersion(version); err != nil {
			return fmt.Sprintf("Restore Game failed: %v", err)
		}
		return "Game restored" + versionSuffix + "."
	default:
		return "Magick!\nSave Game / Restore Game / Save Axil / Restore Axil / Realign Status"
	}
}

// extractVersionLetter finds a real Version-letter token (see save.go's
// doc comment) in an OPTIONS target string: a single-letter word, none
// of which appear among the real keywords this switch already checks
// for (SAVE/RESTORE/GAME/AXIL/REALIGN/STATUS are all multi-letter), so
// no explicit exclusion list is needed. Returns "" if none is present.
func extractVersionLetter(target string) string {
	for word := range strings.FieldsSeq(target) {
		if len(word) == 1 && word[0] >= 'A' && word[0] <= 'Z' {
			return word
		}
	}
	return ""
}
