package world

// Level2Grid is a real, connected 50-cell room graph for the dungeon's
// Level 2, extracted from the same clean, computer-rendered grid map as
// Level1Grid (heavymap-grid-clean.gif — see Level1Grid's doc comment for
// the extraction method and its validation, which applies identically
// here). Cells are addressed the same way: a row letter A-H and a column
// number 1-8, offset into a distinct RoomID range via level2Room so the
// two grids can never collide.
//
// HOW THIS DIFFERS FROM Level1Grid'S STORY: an earlier round tried
// building a Level 2 graph starting from Room of Misery (F4, the
// confirmed real starting room) and found it in a tiny, mostly-isolated
// 7-cell pocket (F3, F4, G3, G4, G5, H4, H5) - which looked at the time
// like Level 2's extraction had failed. Going back and computing EVERY
// connected component in the full 64-cell automated graph (not just the
// one reachable from Room of Misery) showed that conclusion was wrong:
// there's a real, well-connected 50-cell component (78% of the grid,
// comparable to Level1Grid's 44/64) - Room of Misery's small pocket is
// simply not part of it. This file ships that 50-cell component, using
// A1 as an arbitrary stable anchor (not confirmed as any particular
// named room) since it's what New(start) needs, not because A1 is known
// to be special.
//
// NOT included here, each for a specific documented reason:
//   - Room of Misery's own 7-cell pocket (F3/F4/G3/G4/G5/H4/H5) - real,
//     but disconnected from the other 50 cells in this data (same
//     honest-gap pattern as Level1Grid's remaining fragment).
//   - 4 named cells confirmed via text-density measurement (see
//     CLAUDE.md) - Icthys (C3), Flox (D4), Horns (E2), Purity (H3) - and
//     2 more small isolated pockets (C3-C4, C6-D6) - all individually
//     border-checked and confirmed isolated from the 50-cell component,
//     not a detection failure (same pattern as Level1Grid's A7/A8/G5/H5).
//   - Item icon placements: NOT extracted for this file (real follow-up
//     work, not guessed). 3 monster placements WERE added in a later
//     round (Wraith at A5, Slug at C2, Ghost at H6) - each individually
//     verified with a tight per-cell crop against the icon legend, not
//     eyeballed from the full grid view (the same lesson the Room of
//     Claws mistake taught earlier). A pixel-fraction color scan (an
//     icon glyph occupies a small minority of a cell's pixels, unlike a
//     zone's dominant background color or a named-room box's border) was
//     used to find candidates first, then each was confirmed by eye
//     before being trusted - several scan hits turned out to be
//     something else entirely (a room-name box's colored border, a
//     "FIRE!"/"EXIT!" warning label's text color, a stairwell arrow) and
//     were correctly excluded, not just accepted at face value.
//   - Room descriptions: same placeholder convention as everywhere else
//     in this project.
//
// Guards: a later round re-derived this file's own row/column pixel
// calibration directly (validated against A6/A8's own printed labels,
// plus the "AGILE STAIR" box - real, at B8 here, not A7/A8 like
// Level1Grid/Level3Grid) and ran the same monster-color scan used for
// Level 3 across all 64 cells. It exactly reproduced the 3 already-
// shipped monsters (Wraith@A5, Slug@C2, Ghost@H6 all matched precisely -
// good further validation this file's calibration has no Level1Grid-
// style bug), and turned up 5 more candidates. Tight-crop-verified each:
// B1, C8, and D8 are real Guards icons (world.Room.Guards, same
// world.Room.Guards/game.passGuards mechanic added for Level1Grid's D4/
// D7 - see its doc comment for the "GUARDS, DOOR" sourcing); G5 and H5
// are ALSO real Guards icons but sit in Room of Misery's disconnected
// 7-cell pocket (see above), so aren't wired in here to avoid claiming
// reachability this data doesn't support. Two more scan hits (E5, E6)
// tight-crop-verified as a stairwell arrow and "FIRE!" warning text
// respectively, and were correctly excluded, same discipline as the
// monster scan above.
func Level2Grid() *World {
	w := New(level2Room("A1"))
	for _, r := range level2Cells {
		if r.Name == "" {
			r.Name = level2CellCode(r.ID)
		}
		if r.Description == "" {
			r.Description = "(room description not yet extracted from the original)"
		}
		w.AddRoom(r)
	}
	w.Rooms[level2Room("A1")].Visited = true
	return w
}

// level2Room turns a grid cell code like "A1" or "H8" into a stable
// RoomID, offset well clear of both CollodonsPile's and Level1Grid's
// RoomID ranges so none of the three graphs can ever collide.
func level2Room(code string) RoomID {
	row := int(code[0] - 'A')
	col := int(code[1] - '1')
	return RoomID(200000 + row*8 + col)
}

// level2CellCode inverts level2Room's ID encoding, the Level2Grid
// counterpart to Level1Grid's level1CellCode.
func level2CellCode(id RoomID) string {
	offset := int(id) - 200000
	row := offset / 8
	col := offset % 8
	return string(rune('A'+row)) + string(rune('1'+col))
}

var level2Cells = []*Room{
	{ID: level2Room("A1"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("A2"), South: level2Room("B1")}},
	{ID: level2Room("A2"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("A3"), South: level2Room("B2"), West: level2Room("A1")}},
	{ID: level2Room("A3"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("A4"), South: level2Room("B3"), West: level2Room("A2")}},
	{ID: level2Room("A4"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("A5"), South: level2Room("B4"), West: level2Room("A3")}},
	{ID: level2Room("A5"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("A6"), South: level2Room("B5"), West: level2Room("A4")}, Monster: "Wraith", MonsterHealth: 2},
	{ID: level2Room("A6"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("A7"), South: level2Room("B6"), West: level2Room("A5")}},
	{ID: level2Room("A7"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("A8"), South: level2Room("B7"), West: level2Room("A6")}},
	{ID: level2Room("A8"), Level: 2, Exits: map[Direction]RoomID{West: level2Room("A7")}},
	{ID: level2Room("B1"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("B2"), North: level2Room("A1")}, Guards: true},
	{ID: level2Room("B2"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("B3"), North: level2Room("A2"), West: level2Room("B1")}},
	{ID: level2Room("B3"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("B4"), North: level2Room("A3"), West: level2Room("B2")}},
	{ID: level2Room("B4"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("B5"), North: level2Room("A4"), West: level2Room("B3")}},
	{ID: level2Room("B5"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("B6"), North: level2Room("A5"), South: level2Room("C5"), West: level2Room("B4")}},
	{ID: level2Room("B6"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("B7"), North: level2Room("A6"), West: level2Room("B5")}},
	{ID: level2Room("B7"), Level: 2, Exits: map[Direction]RoomID{North: level2Room("A7"), West: level2Room("B6")}},
	{ID: level2Room("C1"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("C2"), South: level2Room("D1")}},
	{ID: level2Room("C2"), Level: 2, Exits: map[Direction]RoomID{South: level2Room("D2"), West: level2Room("C1")}, Monster: "Slug", MonsterHealth: 2},
	{ID: level2Room("C5"), Level: 2, Exits: map[Direction]RoomID{North: level2Room("B5"), South: level2Room("D5")}},
	{ID: level2Room("C7"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("C8"), South: level2Room("D7")}},
	{ID: level2Room("C8"), Level: 2, Exits: map[Direction]RoomID{South: level2Room("D8"), West: level2Room("C7")}, Guards: true},
	{ID: level2Room("D1"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("D2"), North: level2Room("C1"), South: level2Room("E1")}},
	{ID: level2Room("D2"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("D3"), North: level2Room("C2"), West: level2Room("D1")}},
	{ID: level2Room("D3"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("D4"), South: level2Room("E3"), West: level2Room("D2")}},
	{ID: level2Room("D4"), Level: 2, Exits: map[Direction]RoomID{West: level2Room("D3")}},
	{ID: level2Room("D5"), Level: 2, Exits: map[Direction]RoomID{North: level2Room("C5"), South: level2Room("E5")}},
	{ID: level2Room("D7"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("D8"), North: level2Room("C7"), South: level2Room("E7")}},
	{ID: level2Room("D8"), Level: 2, Exits: map[Direction]RoomID{North: level2Room("C8"), West: level2Room("D7")}, Guards: true},
	{ID: level2Room("E1"), Level: 2, Exits: map[Direction]RoomID{North: level2Room("D1"), South: level2Room("F1")}},
	{ID: level2Room("E3"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("E4"), North: level2Room("D3")}},
	{ID: level2Room("E4"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("E5"), West: level2Room("E3")}},
	{ID: level2Room("E5"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("E6"), North: level2Room("D5"), South: level2Room("F5"), West: level2Room("E4")}},
	{ID: level2Room("E6"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("E7"), South: level2Room("F6"), West: level2Room("E5")}},
	{ID: level2Room("E7"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("E8"), North: level2Room("D7"), South: level2Room("F7"), West: level2Room("E6")}},
	{ID: level2Room("E8"), Level: 2, Exits: map[Direction]RoomID{South: level2Room("F8"), West: level2Room("E7")}},
	{ID: level2Room("F1"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("F2"), North: level2Room("E1"), South: level2Room("G1")}},
	{ID: level2Room("F2"), Level: 2, Exits: map[Direction]RoomID{South: level2Room("G2"), West: level2Room("F1")}},
	{ID: level2Room("F5"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("F6"), North: level2Room("E5")}},
	{ID: level2Room("F6"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("F7"), North: level2Room("E6"), South: level2Room("G6"), West: level2Room("F5")}},
	{ID: level2Room("F7"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("F8"), North: level2Room("E7"), South: level2Room("G7"), West: level2Room("F6")}},
	{ID: level2Room("F8"), Level: 2, Exits: map[Direction]RoomID{North: level2Room("E8"), South: level2Room("G8"), West: level2Room("F7")}},
	{ID: level2Room("G1"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("G2"), North: level2Room("F1"), South: level2Room("H1")}},
	{ID: level2Room("G2"), Level: 2, Exits: map[Direction]RoomID{North: level2Room("F2"), South: level2Room("H2"), West: level2Room("G1")}},
	{ID: level2Room("G6"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("G7"), North: level2Room("F6"), South: level2Room("H6")}},
	{ID: level2Room("G7"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("G8"), North: level2Room("F7"), South: level2Room("H7"), West: level2Room("G6")}},
	{ID: level2Room("G8"), Level: 2, Exits: map[Direction]RoomID{North: level2Room("F8"), South: level2Room("H8"), West: level2Room("G7")}},
	{ID: level2Room("H1"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("H2"), North: level2Room("G1")}},
	{ID: level2Room("H2"), Level: 2, Exits: map[Direction]RoomID{North: level2Room("G2"), West: level2Room("H1")}},
	{ID: level2Room("H6"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("H7"), North: level2Room("G6")}, Monster: "Ghost", MonsterHealth: 2},
	{ID: level2Room("H7"), Level: 2, Exits: map[Direction]RoomID{East: level2Room("H8"), North: level2Room("G7"), West: level2Room("H6")}},
	{ID: level2Room("H8"), Level: 2, Exits: map[Direction]RoomID{North: level2Room("G8"), West: level2Room("H7")}},
}
