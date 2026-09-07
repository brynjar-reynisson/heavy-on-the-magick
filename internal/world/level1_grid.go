package world

// Level1Grid is a real, connected 64-cell room graph for the dungeon's
// Level 1, extracted from the clean, computer-rendered grid map
// (heavymap-grid-clean.gif — see known_room_names.go and CLAUDE.md for
// the source and the "room granularity is per-maze-cell" finding that
// prompted this). Cells are addressed the same way the source map
// addresses them: a row letter A-H and a column number 1-8.
//
// HOW CONNECTIVITY WAS EXTRACTED (and its confidence level): the source
// image draws each open passage between adjacent cells as a colored bar
// crossing their shared grid line; a closed border shows no color there
// at all. This was extracted programmatically (sampling pixel colors
// along every shared cell border, not eyeballed cell-by-cell) and the
// method was validated against a connection already confirmed real from
// an independent source: the walkthrough-sourced CollodonsPile room
// "Room of Misery" has a confirmed East exit, and this same map (on
// Level 2, a different grid not included in this file) shows exactly
// that cell (F4) with an open East border to F5 - the extraction method
// correctly found it. That gives real, if not absolute, confidence in
// the Level 1 connections below.
//
// 6 of the 64 cells (A7, A8, F3, G3, G5, H5) showed no detected open
// border via the automated pixel-sampling pass - their distinctively
// colored sub-room border boxes (see the source image) interfered with
// the plain corridor-color detection used for the other 58. Two of them
// (F3 "Room of Stings" and G3 "Exit") were then resolved by careful
// direct visual inspection instead, specifically because the automated
// graph came out split into two disconnected halves at exactly the
// C-row/D-row seam - and F3/G3 sit right at that seam. This matched a
// real, independently-sourced fact rather than contradicting one: the
// CASA walkthrough (see CollodonsPile) already states "Wolfdorp -NW->
// Room of Stings" and "Room of Stings -North-> Morfang", i.e. Room of
// Stings really is the connector between those two areas - good
// cross-confirmation that both the manual read and the two graph halves
// it joins are correct.
//
// A7, A8, G5, and H5 (Agile Stair, Furnace Room, and two unnamed/named
// cells) were later individually re-checked on all 4 sides with the same
// validated pixel-scan method (not re-eyeballed) and confirmed genuinely
// closed on every side - not a detection failure like F3/G3's was. They
// are included as rooms with no Exits, honestly reflecting that no
// ordinary corridor reaches them in this data, not left incomplete out
// of not having looked. A7 (Agile Stair) and A8 (Furnace Room) most
// likely connect via the stairwell/special mechanic shown on the map
// (blue up/down arrow icons appear right at their cell boundaries)
// rather than an ordinary corridor, which this extraction wasn't
// designed to capture. G5 and H5 (Room of Claws) being fully closed on
// every side is more surprising for a named, presumably-reachable room -
// one real, sourced hypothesis worth naming: magic.Demons's Astarot has
// a confirmed "transport to a named location" ability (see
// internal/game's invoke doc comment), so a room reachable only by that
// spell rather than ordinary walking would be thematically consistent
// with what's already confirmed real elsewhere in this project. Not
// confirmed fact, just a plausible explanation for a genuine finding.
//
// Monster placements (Ghost/Werewolf/Vampire/Cyclops) come from the same
// map's icon legend, read precisely against an overlaid coordinate grid
// (not eyeballed) - see CLAUDE.md. MonsterHealth values are placeholders
// (same honesty caveat as CollodonsPile's), not extracted exact values.
//
// CORRECTION: both Werewolf placements were off by one column. A later
// round's tight-crop re-verification (prompted by an unrelated Level 3
// calibration bug turning up the same failure mode elsewhere) found the
// werewolf icons actually sit at C2 and D5, not C3/D6 as first shipped -
// C3 and D6 are genuinely empty cells with their own plain coordinate
// labels visible, confirmed by direct crop. This wasn't a whole-grid
// calibration shift like Level3Grid's (this file's row/column dividers
// and cell labels are independently re-verified correct via A7/A8/F3/H5
// spot checks) - just these 2 monster icons' column got mis-read at
// original extraction time. Fixed once found.
//
// CORRECTION (round 74): "Wraith" renamed to "Vampire" everywhere in
// this project. heavymap-grid-clean.gif's own hand-annotated legend
// (used for every monster placement so far) glossed its red "w" icon as
// "wraith" in plain English, but a second, independent source - a
// composite of REAL in-game screenshots (maps.speccy.cz's "Speccy
// Screenshot Maps" atlas, created by Hippy Smith from actual captured
// gameplay, not hand-drawn) - includes a full "Demons & monsters"
// portrait gallery with the game's own real on-screen creature name
// printed under each portrait. That gallery lists exactly 8 monster
// types; 7 (Troll, Ghost, Slug, Cyclops, Medusa, Werewolf, Wyvern) are
// exact matches for this project's existing roster, and the 8th is
// labeled "VAMPIRE", not "Wraith" - strong evidence this project's
// "Wraith" name came from an approximate fan gloss on the OTHER map's
// legend, not the game's actual name for the creature. (The word
// "WRAITH" does independently appear in the game's real extracted
// 316-word vocabulary too, but so does "VAMPIRE" - both are real
// recognized words; only one has a matching creature portrait, which is
// the deciding evidence here.) Not a new 9th monster - a name correction
// for the same red "w"-icon creature already placed at Level1Grid's F2/
// G1/G2/H1, Level2Grid's A5, Level4Grid's A6, and CollodonsPile's Methos.
//
// "Guards" icons (at D4 and D7, unchanged by the above - they were
// correctly placed originally) are modeled as world.Room.Guards, not
// Monster - guards aren't a BLAST/FREEZE combat encounter. The real
// player interaction is now confirmed too: the game's own in-game hint
// screen (ref_screenshot.png cross-check, see CLAUDE.md) lists
// "GUARDS, DOOR" as an actual example command - wired into game.Handle's
// passGuards, the same TARGET-comma-VERB grammar as "APEX, TALK".
//
// STILL FRAGMENTED: even after the F3/G3 bridge fix, the cell reachable
// from A1 covers 44 of the 64 cells, not all of them - a second group
// (roughly columns 5-8 of rows E through H) remains a separate,
// disconnected component in this data. All 4 originally-unresolved cells
// (A7, A8, G5, H5) have been individually re-verified closed on every
// side, ruling them out as the missing bridge. Going further: EVERY
// possible crossing point between the two components has now been
// individually re-checked by hand (not just the automated pass) - the
// entire column-4/column-5 divide across all 8 rows (A4-A5 through
// H4-H5), and the entire row-D/row-E divide across columns 5-8
// (D5-E5 through D8-E8) - all 12 crossings are confirmed closed. This is
// strong (not just "haven't found it yet") evidence that this fragment
// is genuinely a separate area in this map's own data, not a detection
// failure this project has simply failed to spot - the same way A7/A8's
// isolation turned out to mean "reached via the stairwell mechanic
// instead," this fragment likely connects to the rest of the dungeon via
// something other than an ordinary corridor (another stairwell inside
// the fragment itself is the most likely candidate, given Agile Stair
// was observed appearing at multiple levels' grid boundaries elsewhere
// on this same map). Not chasing this specific bridge further without a
// new source or technique - continuing to re-check ordinary cells by
// hand has diminishing returns once every boundary crossing is ruled out.
//
// NOT wired into CollodonsPile or game.New(): this is a separate,
// standalone graph so it can be shipped without touching the existing
// playable game's tested behavior or RoomID space. Merging it in would
// mean reconciling these 64 cells against CollodonsPile's existing named
// Level-1 rooms (Wolfdorp, Morfang, Pilefoot, Nidus, Room of Stings,
// Room of Arrows), which were sourced from a walkthrough that doesn't
// use this same per-cell addressing - a real follow-up task, not done
// here.
//
// Round 100: re-confirmed via world.SharedNamedRooms (a general cross-
// world scan, not specific to Level 1) that Agile Stair/Room of Stings/
// Room of Arrows remain the ONLY genuinely-named overlaps between this
// file and CollodonsPile - no new candidates turned up. Also pinned down
// the SPECIFIC reason a literal single-graph splice isn't safe yet, not
// just "not done": F3 (Room of Stings) already has a real, pixel-
// extracted North exit to E3 here, but the CASA walkthrough (line ~33
// above) separately states "Room of Stings -North-> Morfang" - two
// independently-sourced sources both claim the SAME direction from the
// SAME named room, but disagree on where it leads. Neither source is
// obviously wrong (E3 is real extracted grid data; Morfang is a real
// walkthrough fact), so forcibly picking one to build a merged graph
// would silently discard the other's real data rather than resolve
// anything - left unresolved rather than guessed, same discipline as
// Sothic Complex's Level 2 vs Level 3 naming clash (see level3_grid.go).
func Level1Grid() *World {
	w := New(level1Room("A1"))
	for _, r := range level1Cells {
		if r.Name == "" {
			r.Name = level1CellCode(r.ID)
		}
		if r.Description == "" {
			r.Description = "(room description not yet extracted from the original)"
		}
		w.AddRoom(r)
	}
	w.Rooms[level1Room("A1")].Visited = true
	return w
}

// level1CellCode inverts level1Room's ID encoding, so unnamed cells can
// still be displayed/mapped by their real grid coordinate (e.g. "A2")
// rather than an empty string.
func level1CellCode(id RoomID) string {
	offset := int(id) - 100000
	row := offset / 8
	col := offset % 8
	return string(rune('A'+row)) + string(rune('1'+col))
}

// level1Room turns a grid cell code like "A1" or "H8" into a stable
// RoomID, offset well clear of CollodonsPile's RoomID range so the two
// graphs can never collide if ever combined.
func level1Room(code string) RoomID {
	row := int(code[0] - 'A')
	col := int(code[1] - '1')
	return RoomID(100000 + row*8 + col)
}

var level1Cells = []*Room{
	{ID: level1Room("A1"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("A2"), South: level1Room("B1")}, Monster: "Ghost", MonsterHealth: 2},
	{ID: level1Room("A2"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("A3"), South: level1Room("B2"), West: level1Room("A1")}},
	{ID: level1Room("A3"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("A4"), South: level1Room("B3"), West: level1Room("A2")}},
	{ID: level1Room("A4"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("A5"), South: level1Room("B4"), West: level1Room("A3")}},
	{ID: level1Room("A5"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("A6"), South: level1Room("B5"), West: level1Room("A4")}},
	{ID: level1Room("A6"), Level: 1, Exits: map[Direction]RoomID{South: level1Room("B6"), West: level1Room("A5")}},
	{ID: level1Room("A7"), Name: "Agile Stair", Level: 1},
	{ID: level1Room("A8"), Name: "Furnace Room", Level: 1},
	{ID: level1Room("B1"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("B2"), North: level1Room("A1"), South: level1Room("C1")}},
	{ID: level1Room("B2"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("B3"), North: level1Room("A2"), South: level1Room("C2"), West: level1Room("B1")}},
	{ID: level1Room("B3"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("B4"), North: level1Room("A3"), South: level1Room("C3"), West: level1Room("B2")}},
	{ID: level1Room("B4"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("B5"), North: level1Room("A4"), South: level1Room("C4"), West: level1Room("B3")}},
	{ID: level1Room("B5"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("B6"), North: level1Room("A5"), South: level1Room("C5"), West: level1Room("B4")}},
	{ID: level1Room("B6"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("B7"), North: level1Room("A6"), South: level1Room("C6"), West: level1Room("B5")}},
	{ID: level1Room("B7"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("B8"), South: level1Room("C7"), West: level1Room("B6")}},
	{ID: level1Room("B8"), Level: 1, Exits: map[Direction]RoomID{South: level1Room("C8"), West: level1Room("B7")}},
	{ID: level1Room("C1"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("C2"), North: level1Room("B1")}},
	{ID: level1Room("C2"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("C3"), North: level1Room("B2"), West: level1Room("C1")}, Monster: "Werewolf", MonsterHealth: 2},
	{ID: level1Room("C3"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("C4"), North: level1Room("B3"), West: level1Room("C2")}},
	{ID: level1Room("C4"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("C5"), North: level1Room("B4"), South: level1Room("D4"), West: level1Room("C3")}},
	{ID: level1Room("C5"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("C6"), North: level1Room("B5"), South: level1Room("D5"), West: level1Room("C4")}},
	{ID: level1Room("C6"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("C7"), North: level1Room("B6"), South: level1Room("D6"), West: level1Room("C5")}, Monster: "Ghost", MonsterHealth: 2},
	{ID: level1Room("C7"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("C8"), North: level1Room("B7"), South: level1Room("D7"), West: level1Room("C6")}},
	{ID: level1Room("C8"), Level: 1, Exits: map[Direction]RoomID{North: level1Room("B8"), South: level1Room("D8"), West: level1Room("C7")}},
	{ID: level1Room("D1"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("D2"), South: level1Room("E1")}},
	{ID: level1Room("D2"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("D3"), South: level1Room("E2"), West: level1Room("D1")}},
	{ID: level1Room("D3"), Level: 1, Exits: map[Direction]RoomID{South: level1Room("E3"), West: level1Room("D2")}},
	{ID: level1Room("D4"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("D5"), North: level1Room("C4"), South: level1Room("E4")}, Guards: true},
	{ID: level1Room("D5"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("D6"), North: level1Room("C5"), West: level1Room("D4")}, Monster: "Werewolf", MonsterHealth: 2},
	{ID: level1Room("D6"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("D7"), North: level1Room("C6"), West: level1Room("D5")}},
	{ID: level1Room("D7"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("D8"), North: level1Room("C7"), West: level1Room("D6")}, Guards: true},
	{ID: level1Room("D8"), Level: 1, Exits: map[Direction]RoomID{North: level1Room("C8"), West: level1Room("D7")}},
	{ID: level1Room("E1"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("E2"), North: level1Room("D1"), South: level1Room("F1")}},
	{ID: level1Room("E2"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("E3"), North: level1Room("D2"), South: level1Room("F2"), West: level1Room("E1")}},
	{ID: level1Room("E3"), Level: 1, Exits: map[Direction]RoomID{North: level1Room("D3"), South: level1Room("F3"), West: level1Room("E2")}},
	{ID: level1Room("E4"), Level: 1, Exits: map[Direction]RoomID{North: level1Room("D4"), South: level1Room("F4")}},
	{ID: level1Room("E5"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("E6"), South: level1Room("F5")}},
	{ID: level1Room("E6"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("E7"), South: level1Room("F6"), West: level1Room("E5")}},
	{ID: level1Room("E7"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("E8"), South: level1Room("F7"), West: level1Room("E6")}},
	{ID: level1Room("E8"), Level: 1, Exits: map[Direction]RoomID{South: level1Room("F8"), West: level1Room("E7")}},
	{ID: level1Room("F1"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("F2"), North: level1Room("E1"), South: level1Room("G1")}},
	{ID: level1Room("F2"), Level: 1, Exits: map[Direction]RoomID{North: level1Room("E2"), South: level1Room("G2"), West: level1Room("F1")}, Monster: "Vampire", MonsterHealth: 2},
	{ID: level1Room("F3"), Name: "Room of Stings", Level: 1, Exits: map[Direction]RoomID{North: level1Room("E3"), East: level1Room("F4"), South: level1Room("G3")}},
	{ID: level1Room("F4"), Level: 1, Exits: map[Direction]RoomID{North: level1Room("E4"), South: level1Room("G4"), West: level1Room("F3")}},
	{ID: level1Room("F5"), Name: "Room of Arrows", Level: 1, Exits: map[Direction]RoomID{North: level1Room("E5")}},
	{ID: level1Room("F6"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("F7"), North: level1Room("E6"), South: level1Room("G6")}},
	{ID: level1Room("F7"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("F8"), North: level1Room("E7"), South: level1Room("G7"), West: level1Room("F6")}},
	{ID: level1Room("F8"), Level: 1, Exits: map[Direction]RoomID{North: level1Room("E8"), South: level1Room("G8"), West: level1Room("F7")}},
	{ID: level1Room("G1"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("G2"), North: level1Room("F1"), South: level1Room("H1")}, Monster: "Vampire", MonsterHealth: 2},
	{ID: level1Room("G2"), Level: 1, Exits: map[Direction]RoomID{North: level1Room("F2"), South: level1Room("H2"), West: level1Room("G1")}, Monster: "Vampire", MonsterHealth: 2},
	{ID: level1Room("G3"), Name: "Exit", Level: 1, Exits: map[Direction]RoomID{North: level1Room("F3"), East: level1Room("G4")}},
	{ID: level1Room("G4"), Level: 1, Exits: map[Direction]RoomID{North: level1Room("F4"), West: level1Room("G3")}},
	{ID: level1Room("G5"), Level: 1},
	{ID: level1Room("G6"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("G7"), North: level1Room("F6"), South: level1Room("H6")}},
	{ID: level1Room("G7"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("G8"), North: level1Room("F7"), South: level1Room("H7"), West: level1Room("G6")}},
	{ID: level1Room("G8"), Level: 1, Exits: map[Direction]RoomID{North: level1Room("F8"), South: level1Room("H8"), West: level1Room("G7")}},
	{ID: level1Room("H1"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("H2"), North: level1Room("G1")}, Monster: "Vampire", MonsterHealth: 2},
	{ID: level1Room("H2"), Level: 1, Exits: map[Direction]RoomID{North: level1Room("G2"), West: level1Room("H1")}},
	{ID: level1Room("H3"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("H4")}},
	{ID: level1Room("H4"), Level: 1, Exits: map[Direction]RoomID{West: level1Room("H3")}},
	{ID: level1Room("H5"), Name: "Room of Claws", Level: 1},
	{ID: level1Room("H6"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("H7"), North: level1Room("G6")}},
	{ID: level1Room("H7"), Level: 1, Exits: map[Direction]RoomID{East: level1Room("H8"), North: level1Room("G7"), West: level1Room("H6")}, Monster: "Cyclops", MonsterHealth: 3},
	{ID: level1Room("H8"), Level: 1, Exits: map[Direction]RoomID{North: level1Room("G8"), West: level1Room("H7")}},
}
