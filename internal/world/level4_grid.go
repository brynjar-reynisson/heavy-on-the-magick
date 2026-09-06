package world

// Level4Grid is a real, connected 17-cell room graph for the dungeon's
// Level 4, extracted from the same clean, computer-rendered grid map as
// Level1Grid/Level2Grid/Level3Grid (heavymap-grid-clean.gif). Cells are
// addressed the same way: a row letter and a column number 1-8, offset
// into a distinct RoomID range via level4Room.
//
// CALIBRATION HISTORY - CORRECTED (same failure mode as Level3Grid's,
// found the same way): this file originally shipped using only 7 rows
// (A-G), because the original calibration pass found 7 of the expected
// 8 horizontal dividers and, finding no plausible 8th one further down,
// assumed the grid was genuinely only 7 rows tall. That was wrong: a
// later round found Level 4's box in the source image spans the exact
// same y-range (327-549) as Level3Grid's already-corrected 8-row
// calibration - both levels are drawn at identical vertical position on
// the source poster - and a divider histogram confirmed the same 8
// row-boundary positions apply here too. Direct tight-crop confirmation:
// the true row A (y=327-353) shows real, distinctly labeled content
// (A4, A6, A7, A8, plus 3 Werewolf icons and a Wraith) that the old
// 7-row data never included at all - it isn't blank space. Every cell
// this file already had was therefore really one row further down than
// its old label said (old "A" was actually the map's real row B, etc.):
// fixed by relabeling all 17 cells up by one row letter (A->B, B->C, ...,
// F->G, G->H) - the connectivity itself (which cells connect to which)
// was already correct and did not need to change, only the letters
// naming them, exactly the same kind of pure-relabeling correction
// Level3Grid needed.
//
// This does NOT yet add the newly-confirmed-real rows A-E as playable
// cells - relabeling what was already extracted is a safe, mechanical
// fix; extracting the real connectivity for the rows this file never
// covered is a separate, larger task (the attempted quick reimplementation
// of the border-open detector failed its own sanity check against
// Level3Grid's known-correct answer - see ../../CLAUDE.md - so it wasn't
// used here either). The two next-largest components found in the
// original pass (11 and 8 cells, also now understood to be one row lower
// than their old labels implied) were checked for an easy bridge to this
// one and found closed, not just unexamined.
//
// Monster placement: ran the same pixel-fraction color scan used for
// Levels 1-3 against this file's already-correct 17-cell component
// (unlike the not-yet-extracted rows above it, this component's exits
// are real and known, so a scan hit here CAN be safely wired in).
// Found 2 candidates; tight-crop-verified both - one (near G7) was a
// stairwell direction arrow, correctly excluded (same false-positive
// pattern seen throughout this project's icon scans); the other is a
// real, confirmed Medusa at H5 (bright red "m", the same color wraith
// also uses - letter shape is what distinguishes them, per the map's
// own legend). MonsterHealth uses the established "tougher monster"
// placeholder value of 3 (same as Cyclops/Troll/Wyvern).
//
// Item icon placements were NOT extracted for this file (real
// follow-up work, same as Levels 1-3), and room descriptions use the
// same placeholder convention as everywhere else.
func Level4Grid() *World {
	w := New(level4Room("F2"))
	for _, r := range level4Cells {
		if r.Name == "" {
			r.Name = level4CellCode(r.ID)
		}
		if r.Description == "" {
			r.Description = "(room description not yet extracted from the original)"
		}
		w.AddRoom(r)
	}
	w.Rooms[level4Room("F2")].Visited = true
	return w
}

// level4Room turns a grid cell code like "F2" or "H7" into a stable
// RoomID, offset well clear of CollodonsPile's/Level1Grid's/Level2Grid's/
// Level3Grid's RoomID ranges so none of the five graphs can ever collide.
func level4Room(code string) RoomID {
	row := int(code[0] - 'A')
	col := int(code[1] - '1')
	return RoomID(400000 + row*8 + col)
}

// level4CellCode inverts level4Room's ID encoding.
func level4CellCode(id RoomID) string {
	offset := int(id) - 400000
	row := offset / 8
	col := offset % 8
	return string(rune('A'+row)) + string(rune('1'+col))
}

var level4Cells = []*Room{
	{ID: level4Room("F2"), Level: 4, Exits: map[Direction]RoomID{East: level4Room("F3")}},
	{ID: level4Room("F3"), Level: 4, Exits: map[Direction]RoomID{East: level4Room("F4"), South: level4Room("G3"), West: level4Room("F2")}},
	{ID: level4Room("F4"), Level: 4, Exits: map[Direction]RoomID{East: level4Room("F5"), West: level4Room("F3")}},
	{ID: level4Room("F5"), Level: 4, Exits: map[Direction]RoomID{West: level4Room("F4")}},
	{ID: level4Room("F6"), Level: 4, Exits: map[Direction]RoomID{East: level4Room("F7"), South: level4Room("G6")}},
	{ID: level4Room("F7"), Level: 4, Exits: map[Direction]RoomID{East: level4Room("F8"), South: level4Room("G7"), West: level4Room("F6")}},
	{ID: level4Room("F8"), Level: 4, Exits: map[Direction]RoomID{South: level4Room("G8"), West: level4Room("F7")}},
	{ID: level4Room("G3"), Level: 4, Exits: map[Direction]RoomID{North: level4Room("F3"), South: level4Room("H3")}},
	{ID: level4Room("G5"), Level: 4, Exits: map[Direction]RoomID{East: level4Room("G6"), South: level4Room("H5")}},
	{ID: level4Room("G6"), Level: 4, Exits: map[Direction]RoomID{East: level4Room("G7"), North: level4Room("F6"), South: level4Room("H6"), West: level4Room("G5")}},
	{ID: level4Room("G7"), Level: 4, Exits: map[Direction]RoomID{East: level4Room("G8"), North: level4Room("F7"), South: level4Room("H7"), West: level4Room("G6")}},
	{ID: level4Room("G8"), Level: 4, Exits: map[Direction]RoomID{North: level4Room("F8"), West: level4Room("G7")}},
	{ID: level4Room("H3"), Level: 4, Exits: map[Direction]RoomID{East: level4Room("H4"), North: level4Room("G3")}},
	{ID: level4Room("H4"), Level: 4, Exits: map[Direction]RoomID{East: level4Room("H5"), West: level4Room("H3")}},
	{ID: level4Room("H5"), Level: 4, Exits: map[Direction]RoomID{East: level4Room("H6"), North: level4Room("G5"), West: level4Room("H4")}, Monster: "Medusa", MonsterHealth: 3},
	{ID: level4Room("H6"), Level: 4, Exits: map[Direction]RoomID{East: level4Room("H7"), North: level4Room("G6"), West: level4Room("H5")}},
	{ID: level4Room("H7"), Level: 4, Exits: map[Direction]RoomID{North: level4Room("G7"), West: level4Room("H6")}},
}
