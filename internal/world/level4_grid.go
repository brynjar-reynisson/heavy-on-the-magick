package world

// Level4Grid is a real, 27-cell room graph for the dungeon's Level 4,
// extracted from the same clean, computer-rendered grid map as
// Level1Grid/Level2Grid/Level3Grid (heavymap-grid-clean.gif). Cells are
// addressed the same way: a row letter and a column number 1-8, offset
// into a distinct RoomID range via level4Room. 17 of the 27 form one
// fully connected component reachable from the start room; the other 10
// (Scales/D2, Doubt of Rabak/D3, a Wyvern/E5, The Crypt/F1, Exit/G2,
// Pride/G4, plus 3 more Wyverns at A1/A3/A5 and a Vampire at A6 - see
// "ROUND 72" below) are real, deliberately isolated special rooms/
// monsters - see the "ROUND 58"/"ROUND 68"/"ROUND 72" sections below.
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
// (A4, A6, A7, A8, plus 3 Werewolf icons and a Vampire) that the old
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
// ROUND 58: applied the "check a gap/newly-confirmed-real area for
// named special rooms" technique (already used on Level2Grid/
// Level3Grid) to the rows above this component. Tight-cropping rows
// A-G found 6 real named special rooms, individually pixel-confirmed:
//   - F4 "The Chasm" is ALREADY a normal cell in the 17-cell main
//     component (real exits East to F5, West to F3) - it was simply
//     missing its Name, the same Flox/D4 situation found in
//     Level2Grid. Fixed by adding the name only, zero connectivity
//     risk.
//   - D2 "Scales", D3 "Doubt of Rabak", F1 "The Crypt", G4 "Pride" are
//     genuinely absent from level4Cells entirely - added as real,
//     named, isolated cells (no Exits - connectivity not extracted,
//     same honest convention as every other special room in this
//     project).
//   - G2 "Exit" - also genuinely absent, added the same way. This is a
//     second real, confirmed Exit location (Level1Grid's G3 was the
//     first) - the manual's confirmed "3 exits" detail means this
//     dungeon has (at least) 2 of them now individually verified real,
//     even though this one, being isolated with no confirmed Exits of
//     its own, can't actually be reached/trigger Game.Won in this
//     file today.
//
// ROUND 68: swept the remaining unchecked rows C-E for anything missed
// by round 58's pass (which had focused on rows A-B and the immediate
// Scales/Doubt of Rabak/Chasm/Exit/Pride area). Found one more real,
// individually pixel-confirmed thing: E5 has a real Wyvern monster icon
// (same confirmed rare exact blue RGB(0,132,255) used elsewhere on this
// map). Added as an isolated cell (no Exits), same convention as the
// round-58 finds. The rest of rows C-E checked in this sweep (C6/C7/C8,
// D5/D6/D7/D8, E6/E8) are plain, unlabeled cells - left unadded.
//
// ROUND 72: zone_monsters.go's independently-sourced ZoneMonsterSightings
// records "Wormring: Wyvern x4" - Wormring being the magenta-bordered
// zone visible spanning roughly row A, columns 1-5 on this same map (see
// known_room_names.go) - but this file only had 1 Wyvern (E5) placed so
// far, since rounds 58/68's monster sweep only covered rows C-E, not
// row A. Re-derived row A's own column boundaries directly (dividers at
// x=373,400,428,455,483,511,539,567,593, reusing the already-validated
// row-A y-band 327-353) and ran the same pixel-fraction color scan
// against it. Found exactly 3 more real, tight-crop-verified Wyverns at
// A1, A3, A5 (blue "w", identical glyph shape to the wyvern legend entry)
// - together with E5 that's 4 total, an exact match for "Wormring:
// Wyvern x4", real cross-validation of both this map reading and that
// independent source. Also found a real Vampire at A6 (red "w" - the
// creature this project called "Wraith" until round 74, see Level1Grid's
// doc comment for the rename - and wyvern share the same lowercase
// glyph shape per the map's own legend, distinguished only by color,
// same as the established vampire/medusa color-sharing case) - A6 sits
// in the yellow "Methos" zone (columns 6-8 of row A), independently
// corroborating last round's zone_monsters.go-sourced Vampire placement
// in CollodonsPile's Methos room from a completely different source
// (this map's own icon data, not the zone-sighting list). A candidate at
// B6 (small red blob) tight-crop-verified as the "up level" arrow icon,
// not a monster - correctly excluded, same false-positive pattern seen
// throughout this project's icon scans. A2 and A4 were also checked and
// are plain cells (A4 shows only its own coordinate label). All 4 new
// finds added as isolated cells (no Exits), same convention as every
// other special room/monster in this file - connectivity for row A
// wasn't extracted.
//
// ROUND 79: zone_monsters.go's independently-sourced ZoneMonsterSightings
// also records "Doubt of Rabak: Vampire x1" (using this project's
// current name for the creature, since round 74's rename) - a real,
// separate entry from "Methos: Vampire x1" (already corroborated by
// A6, round 72) that had sat unaddressed. Doubt of Rabak is already a
// real, named, isolated cell in this file (D3) - added the confirmed
// Vampire to it directly, same convention as every other zone_monsters.go
// cross-reference in this project (Sothic Complex/Ghost and Rook of
// Hydra/Wyvern in Level3Grid, round 71).
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
	{ID: level4Room("F4"), Name: "The Chasm", Level: 4, Exits: map[Direction]RoomID{East: level4Room("F5"), West: level4Room("F3")}},
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

	// 5 real, tight-crop-verified named special rooms in the newly-
	// confirmed-real rows above the 17-cell component (see this file's
	// "round 58" doc update below) - added as isolated cells (no Exits,
	// connectivity not extracted) same as every other special room in
	// this project's grid files.
	{ID: level4Room("D2"), Name: "Scales", Level: 4},
	{ID: level4Room("D3"), Name: "Doubt of Rabak", Level: 4, Monster: "Vampire", MonsterHealth: 2},
	{ID: level4Room("E5"), Level: 4, Monster: "Wyvern", MonsterHealth: 3},
	{ID: level4Room("F1"), Name: "The Crypt", Level: 4},
	{ID: level4Room("G2"), Name: "Exit", Level: 4},
	{ID: level4Room("G4"), Name: "Pride", Level: 4},

	// 4 more real, tight-crop-verified finds in row A (see this file's
	// "round 72" doc update): 3 Wyverns completing the Wormring zone's
	// confirmed count of 4, plus a Vampire in the Methos zone.
	{ID: level4Room("A1"), Level: 4, Monster: "Wyvern", MonsterHealth: 3},
	{ID: level4Room("A3"), Level: 4, Monster: "Wyvern", MonsterHealth: 3},
	{ID: level4Room("A5"), Level: 4, Monster: "Wyvern", MonsterHealth: 3},
	{ID: level4Room("A6"), Level: 4, Monster: "Vampire", MonsterHealth: 2},
}
