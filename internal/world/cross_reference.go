package world

import "regexp"

// cellCodePattern matches the auto-generated fallback name every
// Level1-4Grid gives an unnamed cell (its own row-letter+column-number
// address, e.g. "A1", "H7" — see e.g. Level3Grid's `if r.Name == "" {
// r.Name = level3CellCode(r.ID) }`). All 4 level grids independently use
// this SAME addressing scheme, so "A1" (say) exists as a real Name in
// every one of them — a naming-SCHEME coincidence, not evidence the
// rooms are the same physical place. SharedNamedRooms exists specifically
// to separate that noise from genuine content overlaps (e.g. "Room of
// Stings" appearing in both CollodonsPile and Level1Grid).
var cellCodePattern = regexp.MustCompile(`^[A-H][1-8]$`)

// SharedNamedRooms finds every real (non-cell-code) Room.Name that
// appears in both a and b, returning each match's RoomID in each world.
// A genuinely-named room appearing in two independently-sourced worlds
// is real evidence worth surfacing (this project has repeatedly found
// and acted on exactly this kind of overlap — see e.g. Level1Grid's
// Room of Stings/Room of Arrows/Agile Stair matching CollodonsPile, and
// Level3Grid's Sothic Complex matching CollodonsPile's own), but this
// function deliberately does NOT claim the two RoomIDs are the same
// physical place — see each call site / the doc comments on
// CollodonsPile, Level1Grid, and Level3Grid for the honest, case-by-case
// judgment call already made about each specific overlap (e.g. "Exit"
// appears in both Level1Grid and Level4Grid, but is understood to be two
// DIFFERENT physical rooms — the dungeon has multiple real exits per
// round 19's win-condition find — not a shared location).
func SharedNamedRooms(a, b *World) map[string][2]RoomID {
	shared := map[string][2]RoomID{}
	for idA, roomA := range a.Rooms {
		if roomA.Name == "" || cellCodePattern.MatchString(roomA.Name) {
			continue
		}
		for idB, roomB := range b.Rooms {
			if roomB.Name == roomA.Name {
				shared[roomA.Name] = [2]RoomID{idA, idB}
			}
		}
	}
	return shared
}
