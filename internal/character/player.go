// Package character models Axil the Able and the game's Golden-Dawn-themed
// grade/stat system.
package character

import (
	"math/rand/v2"
	"strings"
)

// Grade is one of the Hermetic Order of the Golden Dawn ranks the game
// borrows for its progression system. Confirmed present as literal strings
// in the game's decompressed memory (see ../../CLAUDE.md, RAM page 5 string
// dump): Neophyte (per contemporary review, the starting grade), Zelator,
// Practicus, Philosophus, Adeptus Minor, Adeptus Major, Adeptus Exemptus,
// Magister Templi, Ipsissimus. The real Golden Dawn system also has a
// "Theoricus" grade between Zelator and Practicus, but that string has not
// been confirmed in the game's data — do not assume it's present until it
// turns up in the disassembly.
//
// OPEN QUESTION (round 91): the real Golden Dawn system also has a
// "Magus" grade between Magister Templi and Ipsissimus - and "MAGUS" IS
// a real, confirmed word in the game's own extracted 316-word parser
// vocabulary (parser.Vocabulary). That alone doesn't confirm it's a
// displayed Grade string though (the vocabulary table is what the
// parser recognizes as INPUT, a different data source than the RAM
// string dump that confirmed the other 9 grade names as OUTPUT text).
// Checked directly: searched all 4 of this repo's .z80 memory snapshots
// for the literal ASCII bytes "MAGUS" - not found in any of them. This
// doesn't disprove it (the game's own text uses a custom, not-yet-fully-
// decoded print engine/encoding - see CLAUDE.md's disassembly notes -
// so a real string could exist without matching a plain-ASCII search),
// but it's not enough to confirm the grade either. Left unresolved and
// NOT added - same "don't assume it's present" discipline as Theoricus.
type Grade int

const (
	Neophyte Grade = iota
	Zelator
	Practicus
	Philosophus
	AdeptusMinor
	AdeptusMajor
	AdeptusExemptus
	MagisterTempli
	Ipsissimus
)

func (g Grade) String() string {
	switch g {
	case Neophyte:
		return "Neophyte"
	case Zelator:
		return "Zelator"
	case Practicus:
		return "Practicus"
	case Philosophus:
		return "Philosophus"
	case AdeptusMinor:
		return "Adeptus Minor"
	case AdeptusMajor:
		return "Adeptus Major"
	case AdeptusExemptus:
		return "Adeptus Exemptus"
	case MagisterTempli:
		return "Magister Templi"
	case Ipsissimus:
		return "Ipsissimus"
	default:
		return "Unknown"
	}
}

// Player is Axil the Able, the game's protagonist.
//
// Skill, Stamina, and Luck are confirmed stat names (found as literal
// strings in the game's memory) and are randomly generated at the start of
// a game per contemporary reviews — confirmed directly by reading a live
// SpecEmu session's stats screen for a fresh Neophyte character (see
// ../../CLAUDE.md): Stamina 36, Skill 8, Luck 4. The official manual
// further confirms these are randomly rolled and can only be re-rolled
// (Option Screen's "Realign Status"), never directly edited. Grade and
// "Experience Points" are also confirmed strings — CORRECTED (round 30):
// this was originally split into two separate fields, Experience and
// Points, but the actual live-observed stats screen shows them as one
// combined label and value ("Experience Points: 0" - see ../../CLAUDE.md),
// not two independent stats, so this is now a single ExperiencePoints
// field matching what's actually confirmed. The exact generation
// formula/range for Skill/Stamina/Luck isn't extracted — NewPlayer
// randomly rolls within estimated ranges centered on that one observed
// sample (see the roll-range constants below), not a fixed canonical
// value and not the real original formula.
type Player struct {
	Name  string
	Grade Grade

	Skill   int
	Stamina int
	// MaxStamina is Stamina's ceiling — confirmed real indirectly: the
	// CASA walkthrough (see ../../CLAUDE.md) instructs casting TRANSFUSION
	// "a few times until you get maximum stamina", implying a real cap
	// exists rather than unlimited healing. Set to the freshly-rolled
	// Stamina value whenever NewPlayer/Realign run, so a fresh/realigned
	// character starts at full health by definition.
	MaxStamina int
	Luck       int

	// ExperiencePoints is the confirmed real "Experience Points" stat
	// (see the correction note above). Not yet awarded by any confirmed
	// original formula — see game.Handle's combat doc comment for how
	// this port currently grants it.
	ExperiencePoints int

	// Items is Axil's carried inventory — real objects can be picked up
	// from a world.Room and carried (see game.Handle's PICKUP/DROP).
	Items []string
}

// HasItem reports whether the player is carrying an item by name
// (case-insensitive, matching how object names are otherwise compared
// throughout this port).
func (p *Player) HasItem(name string) bool {
	for _, item := range p.Items {
		if strings.EqualFold(item, name) {
			return true
		}
	}
	return false
}

// Stat roll ranges. The original randomizes Stamina/Skill/Luck within
// some bounds at game start (confirmed by the manual and by observing a
// live rolled character), but the exact original bounds haven't been
// extracted — these are estimated ranges centered on the one confirmed
// sample (Stamina 36, Skill 8, Luck 4), deliberately kept from rolling
// too low a Stamina (an unplayably-short game) since a real player
// mentioned that's implausible for the original. Not extracted fact,
// just a reasonable placeholder pending the real formula.
const (
	minStamina, maxStamina = 28, 45
	minSkill, maxSkill     = 4, 12
	minLuck, maxLuck       = 1, 8
)

// NewPlayer creates Axil at the start of a game, at the lowest grade, with
// randomly rolled Stamina/Skill/Luck (see the roll-range constants above
// for the honesty caveat on the exact bounds).
func NewPlayer() *Player {
	p := &Player{Name: "Axil", Grade: Neophyte}
	p.Realign()
	return p
}

// Realign rerolls Stamina, Skill, and Luck within the same estimated
// ranges NewPlayer uses. Confirmed real: the official manual's Option
// Screen has a "Realign Status" choice that rerolls these stats (they
// can't be directly edited, only regenerated) — see game.Handle's OPTIONS
// case.
func (p *Player) Realign() {
	p.Stamina = minStamina + rand.IntN(maxStamina-minStamina+1)
	p.MaxStamina = p.Stamina
	p.Skill = minSkill + rand.IntN(maxSkill-minSkill+1)
	p.Luck = minLuck + rand.IntN(maxLuck-minLuck+1)
}

// IsDead reports whether Axil has run out of Stamina — confirmed real:
// the official manual states "If you run out of Stamina, you Die."
func (p *Player) IsDead() bool {
	return p.Stamina <= 0
}
