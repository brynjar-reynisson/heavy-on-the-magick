// Package world models Collodon's Pile: the multi-level dungeon Axil the
// Able is trapped in. Rooms connect via up to 8 compass-direction exits,
// matching the original ZX Spectrum game (per contemporary reviews: "most
// rooms have multiple exits across the 8 compass directions").
//
// This package is also the foundation for the explored-map feature: each
// Room tracks whether the player has visited it, and World can render that
// as a map once enough of the room graph is known. See MapView in map.go.
//
// STATUS: the original game's room/exit data table was never located in
// the disassembly (see ../../CLAUDE.md, "Open next steps"). Real room
// data instead comes from a published walkthrough — see CollodonsPile's
// doc comment for exactly what's confirmed vs. placeholder.
package world

import "maps"

// Direction is one of the 8 compass directions a Room can have an Exit in.
type Direction int

const (
	North Direction = iota
	NorthEast
	East
	SouthEast
	South
	SouthWest
	West
	NorthWest
)

func (d Direction) String() string {
	switch d {
	case North:
		return "North"
	case NorthEast:
		return "NorthEast"
	case East:
		return "East"
	case SouthEast:
		return "SouthEast"
	case South:
		return "South"
	case SouthWest:
		return "SouthWest"
	case West:
		return "West"
	case NorthWest:
		return "NorthWest"
	default:
		return "Unknown"
	}
}

// Opposite returns the direction you'd take to undo a step in d — used both
// for basic movement sanity-checking and for laying out the explored map
// (placing a room one step in the opposite direction of how it was reached).
func (d Direction) Opposite() Direction {
	return (d + 4) % 8
}

// RoomID identifies a single room. The original game almost certainly
// indexes rooms by a small integer (consistent with the picture/object
// table's indexing style found during disassembly); using that same shape
// here keeps the port faithful once the real IDs are known.
type RoomID int

// Room is a single location in the dungeon.
type Room struct {
	ID    RoomID
	Name  string
	Level int // dungeon level, 1 = topmost (per contemporary descriptions of a multi-level dungeon); 0 = unknown/unset
	// Description is honestly still a placeholder everywhere - no room
	// description TEXT has ever been located in the disassembly, despite
	// dozens of rounds' attempts (see CLAUDE.md's "Open next steps": the
	// string-print loop and movement/room dispatcher are still untraced).
	// Worth naming a real, growing body of circumstantial evidence
	// (round 95) that this may not be a gap so much as a wrong
	// assumption: the game is classified "Adventure/Graphics" (round 74's
	// web search), the confirmed real screenshot atlas
	// (heavymap-speccy-screenshots.png) shows each room as a drawn
	// corridor SCENE with no on-screen text area at all, and a targeted
	// high-bit-masked memory search for known real room/zone names
	// (round 95) found them all sitting inside the already-known 316-
	// word parser vocabulary table, not a separate prose/message table -
	// consistent with "the room's real content IS the picture," not with
	// a hidden text table nobody's found yet. Not confirmed either way -
	// still an honest placeholder, not rewritten to claim more than is
	// known - but if true, the real remaining "faithful room content"
	// gap is graphics fidelity (see internal/graphics's ongoing work),
	// not undiscovered prose.
	Description string
	Exits       map[Direction]RoomID

	// Visited is true once the player has entered this room. Drives both
	// gameplay (e.g. "MONSTER NEARBY" logic seen in the disassembly may
	// depend on room state) and the explored-map feature.
	Visited bool

	// DoorPasswords are words that unlock progress when said via
	// "DOOR, <password>" in this room — real, confirmed in-game
	// interactions sourced from the CASA walkthrough (see
	// collodons_pile.go and CLAUDE.md), not invented. Empty for most
	// rooms (no known password puzzle, or none exists there).
	DoorPasswords []string

	// DoorHints are the real riddle/clue text Apex (or a room's "guard"
	// pillars) gives for this room's real DoorPasswords, surfaced via
	// "APEX, DOOR" (round 159; see game.apexDoorHint) — one of the 3
	// confirmed example commands from the game's own real hint screen
	// (already quoted verbatim in game.help(), but never actually wired
	// to a response until now). Source: The CRPG Addict's 2016 blog
	// post, quoting the actual riddles the game presented and their
	// real solutions — "CRY AND ENTER DOOR" (answer WOLF, a "cry wolf"
	// pun), "TO ENTER IS MADNESS" (answer LUNACY), and "TO ENTER SAY A
	// NUMBER OF MAGICK WORDS" (answer ELEVEN, tying directly back to
	// the manual's own separately-confirmed "the number of Magick is
	// 11" fact — round 110 — which had no confirmed mechanical use
	// until this cross-reference). The riddle CONTENT is real and
	// sourced; the exact on-screen casing/punctuation is NOT confirmed
	// to the same pixel-exact standard as game.help()'s screen text
	// (that was cross-checked against a real screenshot; this is a
	// blogger's own prose quoting the game from memory/notes), so
	// stored here in normal sentence case rather than claimed ALL-CAPS
	// verbatim. Empty for rooms with no confirmed riddle text (most
	// DoorPasswords rooms, including Secunda Porta's — no riddle for it
	// has been found by any source checked so far).
	DoorHints []string

	// TollItem is the item name this room's door requires dropped/placed
	// to open - real, distinct confirmed mechanic, separate from
	// DoorPasswords (which are typed words, not carried items). Two
	// independent sources agree on the mechanic: Spectrum Computing's
	// instructions file ("For a door with a toll sign by it (ask apex) a
	// bag of gold is the key (put it on the table)") and a fresh, more
	// detailed CASA walkthrough re-read (round 64) that found the real
	// trigger phrase used repeatedly is literally "DROP <item>"
	// ("EXAMINE TABLE, DROP KEY (door opens)", also seen with BAG and
	// SLAT at other rooms) - i.e. this isn't only a "Toll sign" special
	// case, it's the same general item-drop mechanic recurring at
	// several rooms with no toll sign mentioned at all. See
	// game.drop/game.payToll for where this is checked. Empty for rooms
	// with no known required-drop door.
	TollItem string

	// Monster is the name of a creature present in this room, and
	// MonsterHealth how many BLAST hits it takes to defeat it — real,
	// sourced from the CASA walkthrough (confirmed as "BLAST (until
	// monster dies)", implying multiple hits, not a one-shot kill).
	// MonsterHealth's exact original value isn't stated in the source, so
	// a small placeholder count is used where a monster is confirmed
	// present but the real hit count isn't known — see collodons_pile.go.
	// Empty Monster means no known creature in this room.
	Monster       string
	MonsterHealth int

	// Items are pickupable object names present in this room — real,
	// sourced data (see collodons_pile.go for provenance), removed from
	// here and added to the player's inventory on a successful PICKUP.
	Items []string

	// Guards, when true, means this room has a guards obstacle - a real
	// icon on the clean grid map's own legend (a red "guards" glyph,
	// distinct from any monster letter). The real player interaction is
	// confirmed too, from the game's own in-game hint screen (see
	// ../../CLAUDE.md's ref_screenshot.png cross-check): "GUARDS, DOOR"
	// is one of its listed example commands, the same TARGET-comma-VERB
	// grammar as "APEX, TALK". What exactly happens beyond "you get past"
	// isn't stated (no payment/password precondition confirmed), so
	// game.Handle's guards() models the simplest honest reading: saying
	// it clears them. Set for real, tight-crop-verified placements only
	// (see level1_grid.go).
	Guards bool

	// Fire, when true, means this room is a real fire hazard blocking
	// ordinary movement - two independent sources agree it's real: the
	// CASA walkthrough states, of the Clasp (a real, already-placed item
	// in Trollwynd), "Pick up CLASP (this enables you to walk through
	// fire)"; separately, a tight-crop of the clean grid map's Level 2
	// section confirms a real "FIRE!" warning label at cell D6 (round 80
	// — a prior round's doc comment had mis-described this as E6, an
	// off-by-one-row description error, not a data bug, since no Fire
	// field existed yet to be wrong; corrected when this field was
	// added). No source states what happens to a player without the
	// Clasp beyond "enables you to walk through" - game.move models the
	// simplest honest reading (blocks passage), the same convention
	// already used for Guards' "simplest honest reading" precedent. A
	// THIRD independent source cross-confirms fire is real dungeon
	// content (round 84): level_items.go's LevelOneItems (from a
	// completely different map, heavymap-levels1-2.jpg) includes a plain
	// "Fire" entry - on Level 1, not Level 2, so it's not the same D6
	// cell specifically (this project's sources have disagreed on exact
	// Level numbers before, e.g. Room of Misery/Sothic Complex), but it's
	// good independent evidence this hazard is a real, recurring dungeon
	// feature, not a one-off.
	Fire bool

	// HasTable is true for a real, sourced room fixture: the CASA
	// walkthrough repeatedly uses "EXAMINE TABLE" as a command in
	// specific named rooms (Room of Misery, Trollwynd, Sothic Complex,
	// Methos, Wolfdorp, Room of Stings, Morfang, Room of Arrows — see
	// collodons_pile.go), always right before either picking up an item
	// or dropping one to pay a real TollItem ("put it on the table" —
	// the same table, matching the instructions file's own phrasing for
	// that mechanic). Not an Item (can't be picked up), just a real
	// environmental fact confirmed for these specific rooms - the
	// walkthrough happens to use this phrase at the start of nearly
	// every room visited, so other rooms may well have one too, but
	// only these are actually confirmed; not assumed for the rest.
	HasTable bool

	// HasChest is the same kind of real, sourced fixture as HasTable, for
	// a second distinct container the CASA walkthrough confirms: "EXAMINE
	// CHEST" appears before picking up Garlic in Wolfdorp and Slat in
	// Morfang (round 78) - a real, separate phrase from "EXAMINE TABLE",
	// not a synonym for it. Same honesty convention as HasTable: only set
	// for these 2 specifically confirmed rooms, not assumed elsewhere.
	HasChest bool

	// SwapItem and RevealItem model a real, sourced "protected item"
	// mechanic (round 132): World of Spectrum's separate plain-text
	// instructions file (see game.go's checkSwapItem doc comment for the
	// full sourcing) states, as tips, "To get the Pellet swap it for a
	// Ball" and "Get the Shell and swap it for the egg" - and the same
	// file's own walkthrough shows the identical pattern for Nougat/
	// Nugget (DROP NOUGAT, then a separate PICK UP NUGGET). The fan-made
	// numbered map poster's own key list independently corroborates all
	// 3 exact item names, each qualified "protected" ("Egg - rock,
	// protected", "Pellet - rock, protected", "Nugget (silver), rock,
	// protected" - internal/world/numbered_room_contents.go #12/#31/#49)
	// - real, cross-source confirmation this is a genuine mechanic, not
	// a walkthrough author's turn of phrase. Dropping SwapItem in a room
	// that has RevealItem set reveals it (adds it to Items) - see
	// game.checkSwapItem. Neither field is set on any real, currently-
	// shipped room: the numbered map's cell numbers for this triad
	// (#12/#31/#49) haven't been cross-referenced to a specific
	// CollodonsPile/LevelNGrid room yet (real, scoped follow-up work,
	// same "mechanic real, not yet reachable" pattern already used for
	// TollItem/Fire/Guards before their first real placement).
	SwapItem   string
	RevealItem string
}

// clone returns a deep copy of r, safe to mutate independently of the
// original — see AddRoom's doc comment for why this matters.
func (r *Room) clone() *Room {
	c := *r
	if r.Exits != nil {
		c.Exits = make(map[Direction]RoomID, len(r.Exits))
		maps.Copy(c.Exits, r.Exits)
	}
	if r.DoorPasswords != nil {
		c.DoorPasswords = append([]string(nil), r.DoorPasswords...)
	}
	if r.DoorHints != nil {
		c.DoorHints = append([]string(nil), r.DoorHints...)
	}
	if r.Items != nil {
		c.Items = append([]string(nil), r.Items...)
	}
	return &c
}
