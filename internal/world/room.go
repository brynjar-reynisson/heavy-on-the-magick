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
	ID          RoomID
	Name        string
	Level       int // dungeon level, 1 = topmost (per contemporary descriptions of a multi-level dungeon); 0 = unknown/unset
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
	if r.Items != nil {
		c.Items = append([]string(nil), r.Items...)
	}
	return &c
}
