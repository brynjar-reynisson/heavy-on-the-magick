package world

// Level3Grid is a real, 44-cell room graph for the dungeon's Level 3,
// extracted from the same clean, computer-rendered grid map as
// Level1Grid/Level2Grid (heavymap-grid-clean.gif). Cells are addressed
// the same way: a row letter A-H and a column number 1-8, offset into a
// distinct RoomID range via level3Room. 41 of the 44 form one fully
// connected component reachable from the start room; the other 3 (D4
// "Sothic Complex", F3 "Nani", F5 "Hydra") are real, named, deliberately
// isolated special rooms - see the "SOTHIC COMPLEX" section below and
// the round-57 update.
//
// ROUND 57: applying the same "check for a named special room explaining
// a gap" technique that found Sothic Complex, tight-cropped the F2-F6
// span - entirely absent from the main component, same as D4 was -
// and found 2 more real special rooms: F3 reads "NANI" (Room of Nani,
// a zone name already visible elsewhere on this map) and F5 reads
// "HYDRA" (Rook of Hydra, likewise). F2, F4, and F6 were also checked
// but show no name (F4 just its own plain "F4" coordinate code; F2 has
// a special-room-style border but no legible text in this crop; F6 is
// a plain cell) - left unadded rather than guess at a name or force
// connectivity for them.
//
// CALIBRATION HISTORY - CORRECTED (this matters for anyone diffing old
// output against this file): an earlier round found only 7 of the
// expected 8 horizontal grid dividers and, after searching (wrongly) for
// a missing 8th divider below the others, shipped a smaller 34-cell
// component using a guessed, unusually-tall final row. That guess was
// wrong in a specific, now-understood way: the missing divider wasn't at
// the bottom, it was at the TOP - the grid's real row A starts about one
// row-height above where the earlier calibration began (confirmed by
// finding the actual "A1"/"A2"/"A5"/"A7" labels one row higher than
// assumed, right where the previous calibration's first row was
// mislabeling the map's real row B as "A"). With the true 8-row
// calibration, computing every connected component again found a
// larger, real 41-cell component (73% of the grid, better than the
// earlier miscalibrated 34) - this file replaces the old data entirely
// with the corrected version. The earlier 34-cell version was internally
// consistent (a valid connected maze) but every cell's letter was
// effectively shifted by one row relative to the source map's own
// labeling - worth knowing if anything ever needs to be cross-referenced
// against an old snapshot of this file.
//
// Monster placements: found via the same pixel-fraction color scan used
// for Level2Grid, cross-checked against the map's own printed legend
// (troll/cyclops share one dark-olive RGB(132,132,0), wraith/medusa
// share one bright-red RGB(255,0,0), ghost is bright green RGB(0,255,0)
// vs. slug's dark green RGB(0,132,0) - same-color pairs are only
// distinguishable by the letter GLYPH, not color, so every candidate
// was individually tight-cropped and visually read before being
// trusted). One color-match candidate (a red pixel cluster near D7)
// tight-crop-verified as a stairwell direction arrow, not a monster,
// and was correctly excluded - the same discipline that caught false
// positives in Level2Grid's scan. 7 held up: Wyvern at B1, Ghost at B2
// and E2, Troll at C4/C6/E7/F8. MonsterHealth values follow the same
// placeholder convention as elsewhere in this port (the exact original
// numbers aren't stated in any source found so far) - Troll and Wyvern
// use Cyclops's established "tougher monster" value of 3, Ghost uses
// the common value of 2.
//
// Item placement: cross-referenced the fan-made NUMBERED map/key poster
// (internal/world/numbered_room_contents.go's source - see its doc
// comment) against this file's zone names for the first time. Its
// numbered cell #32 ("Cabinet (Mantis)") sits, visually confirmed by
// tight-cropping the poster's own hand-drawn Level 3 grid, within the
// "GORBURG" zone label's own cell cluster - the same zone this file's
// A1/A2/B1/B2/C1/C2/C3/D1/E1 cells belong to (per the clean grid map's
// magenta zone coloring). "Mantis" is Belezbar's confirmed real Charm
// (see magic.Demons) - this is a genuine, sourced item placement, not
// invented, but it's only ZONE-level confidence, not exact-cell
// confidence (the numbered poster's hand-drawn grid doesn't align
// cleanly enough with this file's lettered cells to say which exact
// cell #32 is - same honest precision limit CollodonsPile's own
// zone-abstracted rooms already accept). Placed on A1, the zone's first
// cell, as a real, working Item - this is the first Charm placement in
// the whole port, making Belezbar's already-implemented (game.invoke)
// but previously-unreachable-in-practice invocation actually succeed
// for a player who finds and carries it. The numbered map also names a
// second item in the same Gorburg zone (#21, "Cabinet (clasp -
// Salamander charm)") that does NOT match any of the 4 confirmed
// demons' Charms (Erlstone/Sword/Mantis/Sunflower) - a real, sourced
// fact left unplaced rather than force a guess at what it's for.
//
// SOTHIC COMPLEX (D4): the map itself explains why D4 was never found
// by the automated border-color-span extraction that produced the
// other 41 cells (C4 has no South exit, E4 has no North exit - the
// "gap" they'd normally close). Directly tight-cropping that grid
// position shows a "SOTHIC COMPLEX" label drawn as an irregular special
// room shape, not a standard grid box - the same reason Level1Grid
// needed manual handling for Agile Stair/Furnace Room/Room of Stings/
// Exit. Added here as a real, named room with no Exits (connectivity
// genuinely not confirmed - the visible diagonal/directional arrows
// near it are the same ambiguous stairwell-style markers this project
// has consistently declined to interpret as ordinary corridors
// elsewhere, e.g. Level1Grid's A7/A8 doc comment).
//
// REAL CROSS-SOURCE DISCREPANCY, left unresolved rather than guessed
// at: CollodonsPile (sourced independently from the CASA walkthrough)
// already has its own "Sothic Complex" room at Level 2, and this clean
// map's own Level 2 grid section separately shows "Sothic Complex" as
// a whole named zone there too - so the name is confirmed real on BOTH
// Level 2 and Level 3 of this same map. Whether that's the same
// physical location (like Agile Stair, confirmed elsewhere to span
// multiple levels via stairwells) or two distinct rooms that happen to
// share a name isn't something any source found so far settles, so
// both are kept as-is rather than "resolving" the discrepancy by
// guessing.
//
// Same honest scope as Level2Grid otherwise: room descriptions use the
// same placeholder convention as everywhere else, and no other item
// placements were extracted for this file.
func Level3Grid() *World {
	w := New(level3Room("A1"))
	for _, r := range level3Cells {
		if r.Name == "" {
			r.Name = level3CellCode(r.ID)
		}
		if r.Description == "" {
			r.Description = "(room description not yet extracted from the original)"
		}
		w.AddRoom(r)
	}
	w.Rooms[level3Room("A1")].Visited = true
	return w
}

// level3Room turns a grid cell code like "A1" or "F8" into a stable
// RoomID, offset well clear of CollodonsPile's/Level1Grid's/Level2Grid's
// RoomID ranges so none of the four graphs can ever collide.
func level3Room(code string) RoomID {
	row := int(code[0] - 'A')
	col := int(code[1] - '1')
	return RoomID(300000 + row*8 + col)
}

// level3CellCode inverts level3Room's ID encoding.
func level3CellCode(id RoomID) string {
	offset := int(id) - 300000
	row := offset / 8
	col := offset % 8
	return string(rune('A'+row)) + string(rune('1'+col))
}

var level3Cells = []*Room{
	{ID: level3Room("A1"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("A2"), South: level3Room("B1")}, Items: []string{"Mantis"}},
	{ID: level3Room("A2"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("A3"), South: level3Room("B2"), West: level3Room("A1")}},
	{ID: level3Room("A3"), Level: 3, Exits: map[Direction]RoomID{South: level3Room("B3"), West: level3Room("A2")}},
	{ID: level3Room("A4"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("A5"), South: level3Room("B4")}},
	{ID: level3Room("A5"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("A6"), South: level3Room("B5"), West: level3Room("A4")}},
	{ID: level3Room("A6"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("A7"), South: level3Room("B6"), West: level3Room("A5")}},
	{ID: level3Room("A7"), Level: 3, Exits: map[Direction]RoomID{South: level3Room("B7"), West: level3Room("A6")}},
	{ID: level3Room("B1"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("B2"), North: level3Room("A1"), South: level3Room("C1")}, Monster: "Wyvern", MonsterHealth: 3},
	{ID: level3Room("B2"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("B3"), North: level3Room("A2"), South: level3Room("C2"), West: level3Room("B1")}, Monster: "Ghost", MonsterHealth: 2},
	{ID: level3Room("B3"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("B4"), North: level3Room("A3"), South: level3Room("C3"), West: level3Room("B2")}},
	{ID: level3Room("B4"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("B5"), North: level3Room("A4"), South: level3Room("C4"), West: level3Room("B3")}},
	{ID: level3Room("B5"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("B6"), North: level3Room("A5"), South: level3Room("C5"), West: level3Room("B4")}},
	{ID: level3Room("B6"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("B7"), North: level3Room("A6"), South: level3Room("C6"), West: level3Room("B5")}},
	{ID: level3Room("B7"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("B8"), North: level3Room("A7"), South: level3Room("C7"), West: level3Room("B6")}},
	{ID: level3Room("B8"), Level: 3, Exits: map[Direction]RoomID{South: level3Room("C8"), West: level3Room("B7")}},
	{ID: level3Room("C1"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("C2"), North: level3Room("B1"), South: level3Room("D1")}},
	{ID: level3Room("C2"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("C3"), North: level3Room("B2"), South: level3Room("D2"), West: level3Room("C1")}},
	{ID: level3Room("C3"), Level: 3, Exits: map[Direction]RoomID{North: level3Room("B3"), South: level3Room("D3"), West: level3Room("C2")}},
	{ID: level3Room("C4"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("C5"), North: level3Room("B4")}, Monster: "Troll", MonsterHealth: 3},
	{ID: level3Room("C5"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("C6"), North: level3Room("B5"), South: level3Room("D5"), West: level3Room("C4")}},
	{ID: level3Room("D4"), Name: "Sothic Complex", Level: 3},
	{ID: level3Room("C6"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("C7"), North: level3Room("B6"), South: level3Room("D6"), West: level3Room("C5")}, Monster: "Troll", MonsterHealth: 3},
	{ID: level3Room("C7"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("C8"), North: level3Room("B7"), South: level3Room("D7"), West: level3Room("C6")}},
	{ID: level3Room("C8"), Level: 3, Exits: map[Direction]RoomID{North: level3Room("B8"), South: level3Room("D8"), West: level3Room("C7")}},
	{ID: level3Room("D1"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("D2"), North: level3Room("C1"), South: level3Room("E1")}},
	{ID: level3Room("D2"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("D3"), North: level3Room("C2"), South: level3Room("E2"), West: level3Room("D1")}},
	{ID: level3Room("D3"), Level: 3, Exits: map[Direction]RoomID{North: level3Room("C3"), South: level3Room("E3"), West: level3Room("D2")}},
	{ID: level3Room("D5"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("D6"), North: level3Room("C5"), South: level3Room("E5")}},
	{ID: level3Room("D6"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("D7"), North: level3Room("C6"), South: level3Room("E6"), West: level3Room("D5")}},
	{ID: level3Room("D7"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("D8"), North: level3Room("C7"), South: level3Room("E7"), West: level3Room("D6")}},
	{ID: level3Room("D8"), Level: 3, Exits: map[Direction]RoomID{North: level3Room("C8"), South: level3Room("E8"), West: level3Room("D7")}},
	{ID: level3Room("E1"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("E2"), North: level3Room("D1"), South: level3Room("F1")}},
	{ID: level3Room("E2"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("E3"), North: level3Room("D2"), West: level3Room("E1")}, Monster: "Ghost", MonsterHealth: 2},
	{ID: level3Room("E3"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("E4"), North: level3Room("D3"), West: level3Room("E2")}},
	{ID: level3Room("E4"), Level: 3, Exits: map[Direction]RoomID{West: level3Room("E3")}},
	{ID: level3Room("E5"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("E6"), North: level3Room("D5")}},
	{ID: level3Room("E6"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("E7"), North: level3Room("D6"), West: level3Room("E5")}},
	{ID: level3Room("E7"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("E8"), North: level3Room("D7"), South: level3Room("F7"), West: level3Room("E6")}, Monster: "Troll", MonsterHealth: 3},
	{ID: level3Room("E8"), Level: 3, Exits: map[Direction]RoomID{North: level3Room("D8"), South: level3Room("F8"), West: level3Room("E7")}},
	{ID: level3Room("F1"), Level: 3, Exits: map[Direction]RoomID{North: level3Room("E1")}},
	{ID: level3Room("F3"), Name: "Nani", Level: 3},
	{ID: level3Room("F5"), Name: "Hydra", Level: 3},
	{ID: level3Room("F7"), Level: 3, Exits: map[Direction]RoomID{East: level3Room("F8"), North: level3Room("E7")}},
	{ID: level3Room("F8"), Level: 3, Exits: map[Direction]RoomID{North: level3Room("E8"), West: level3Room("F7")}, Monster: "Troll", MonsterHealth: 3},
}
