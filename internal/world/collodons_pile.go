package world

// CollodonsPile is the real dungeon from the original game — Axil is
// magically banished to the dungeons beneath "a dreary castle called
// Collodon's Pile" (per contemporary descriptions of the game).
//
// SOURCE AND HONESTY NOTE: unlike everything else extracted for this
// project, this room data does NOT come from disassembling the game's
// Z80 code/memory (the real room/exit table has not been located there
// despite repeated attempts — see ../../CLAUDE.md). It comes from a
// published walkthrough at the Classic Adventures Solution Archive
// (solutionarchive.com, "Heavy_on_the_Magick.txt"), which documents one
// specific playthrough's path through the dungeon. Room *names* and the
// *directions* used to move between them are facts extracted from that
// walkthrough (not a copy of its prose) and cross-confirmed against room
// names we separately found as literal strings in the game's own memory
// during disassembly (e.g. "TROLLWYND", "SOTHIC", "WOLFDORP" all appear
// in both places) — so these are real, correct room names, not invented.
// (Two of these names were corrected against that same vocabulary list
// in round 15 — see the note by Nidus, below.)
//
// SCALE: Hardcore Gaming 101's coverage of this game states it has **255
// distinct rooms** — so these 13 are a small, honest fraction (~5%) of the
// full dungeon, not most of it. "Dozens" (a contemporary review's phrasing)
// undersold it. Extracting the rest would need either a more exhaustive
// walkthrough than the one used here (this one traces a single efficient
// solution path, not every room) or finally locating the real room table
// in the disassembly.
//
// What this data does NOT capture:
//   - Only 13 of 255 rooms (see above).
//   - Only directions actually stated in the walkthrough are recorded (a
//     second, more thorough extraction pass added 2 more: Trollwynd's
//     South exit and Room of Arrows' North exit — both directly stated,
//     not inferred). Two further transitions the second pass surfaced
//     were explicitly hedged by the extraction itself ("or nearby
//     connection", "implied through backtracking") and are deliberately
//     NOT included — they'd just be guessing the reverse of an existing
//     exit, which these games are confirmed (see CLAUDE.md) to sometimes
//     deliberately NOT have (non-Euclidean/one-way shortcuts), so a
//     reverse exit is never assumed automatically, only added when a
//     source directly states it.
//   - Room descriptions (the actual text shown on screen) are NOT sourced
//     from the walkthrough (that would mean copying its prose) and have
//     not been separately extracted from the disassembly yet, so
//     Description is left as an honest placeholder here, not invented
//     flavor text.
//
// DoorPasswords ARE real, sourced verb-resolution data (a third extraction
// pass over the same walkthrough, asked specifically for "TARGET, VERB"
// commands and their stated effects): Secunda Porta's door responds to
// "DOOR, SILENCE"; Wolfdorp's to "DOOR, WOLF" or "DOOR, LUNACY" (the
// walkthrough uses both at different points — modeled as either working,
// since it isn't clear from the extraction whether they're alternates or
// a two-step puzzle); Pilefoot's to "DOOR, ELEVEN". See internal/game for
// where these get checked. The same extraction also confirmed BLAST
// (combat) and TRANSFUSION (restores Stamina) as real standalone verbs
// used repeatedly throughout — see game.Handle.
//
// Monster/MonsterHealth are likewise real per a fourth extraction pass
// (asked for monster encounters + which room + defeat method): a
// (unnamed in the source) monster in Trollwynd and a named Cyclops in
// Nidus, both confirmed "BLAST (until monster dies)" — i.e. multiple
// hits, not one-shot. The exact original hit count isn't stated, so
// MonsterHealth=3 here is a placeholder guess, not extracted fact.
//
// CORRECTED (round 15): this room was originally ported as "Midus" and
// the Level 2 room below as "Solthic Complex", both transcribed from the
// walkthrough. The game's own 316-word vocabulary (extracted directly
// from memory during disassembly, independent of the walkthrough or any
// map — see ../../CLAUDE.md) contains "NIDUS" and "SOTHIC", and does NOT
// contain "MIDUS" or "SOLTHIC" at all. Since the vocabulary is extracted
// straight from the game's own data, it's more authoritative than a
// walkthrough transcription, so both names are corrected here to match
// it: Nidus, Sothic Complex.
//
// Items are real per a fifth extraction pass (asked for pickupable
// objects + which room): Grimoire in Room of Misery, Garlic/Bag/Loaf in
// Wolfdorp, Slat in Morfang, Nugget in Methos, Clasp in Trollwynd. A few
// other items the same pass surfaced ("Scroll", "Nougat", "Key") were
// reported with vague/multiple locations ("two instances", "multiple
// locations") rather than one specific room, so they're deliberately NOT
// placed here — assigning them to a guessed room would be fabrication.
//
// A fan-made numbered map poster (see numbered_room_contents.go) provides
// good independent cross-confirmation and one open discrepancy worth
// noting: it independently confirms Room of Misery as the dungeon's
// START room, on a grid section headed "LEVEL 2" (matching Level: 2
// here); but it also places Agile Stair within a "LEVEL 1"-headed grid
// section, while this file has Agile Stair at Level 4. That poster
// labels rooms as spanning grid-section boundaries in a few places
// (consistent with Agile Stair being a stairwell connecting multiple
// levels, per its own "Stair Wells: Level to Level" legend), so this
// isn't treated as a confirmed correction — Agile Stair's Level here is
// left unchanged rather than guessed which source is more authoritative.
// The poster also spells this room "Sothic Complex" — a third source
// agreeing with the vocabulary-based correction above.
func CollodonsPile() *World {
	w := New(roomMisery)
	for _, r := range []*Room{
		{ID: roomMisery, Name: "Room of Misery", Level: 2, Exits: map[Direction]RoomID{East: roomSecundaPorta}, Items: []string{"Grimoire"}},
		{ID: roomSecundaPorta, Name: "Secunda Porta", Level: 2, Exits: map[Direction]RoomID{North: roomTrollwynd}, DoorPasswords: []string{"SILENCE"}},
		{ID: roomTrollwynd, Name: "Trollwynd", Level: 3, Exits: map[Direction]RoomID{North: roomAgileStair, South: roomSothicComplex}, Monster: "monster", MonsterHealth: 3, Items: []string{"Clasp"}},
		{ID: roomAgileStair, Name: "Agile Stair", Level: 4, Exits: map[Direction]RoomID{SouthEast: roomMethos}},
		{ID: roomMethos, Name: "Methos", Level: 4, Exits: map[Direction]RoomID{South: roomSothicComplex}, Items: []string{"Nugget"}},
		{ID: roomSothicComplex, Name: "Sothic Complex", Level: 2, Exits: map[Direction]RoomID{South: roomWolfdorp}},
		{ID: roomWolfdorp, Name: "Wolfdorp", Level: 1, Exits: map[Direction]RoomID{NorthWest: roomStings}, DoorPasswords: []string{"WOLF", "LUNACY"}, Items: []string{"Garlic", "Bag", "Loaf"}},
		{ID: roomStings, Name: "Room of Stings", Level: 1, Exits: map[Direction]RoomID{North: roomMorfang}},
		{ID: roomMorfang, Name: "Morfang", Level: 1, Exits: map[Direction]RoomID{East: roomArrows}, Items: []string{"Slat"}},
		{ID: roomArrows, Name: "Room of Arrows", Level: 1, Exits: map[Direction]RoomID{East: roomNidus, North: roomWolfdorp}},
		{ID: roomNidus, Name: "Nidus", Level: 1, Exits: map[Direction]RoomID{West: roomPilefoot}, Monster: "Cyclops", MonsterHealth: 3},
		{ID: roomPilefoot, Name: "Pilefoot", Level: 1, Exits: map[Direction]RoomID{North: roomPileCollodom}, DoorPasswords: []string{"ELEVEN"}},
		{ID: roomPileCollodom, Name: "Pile Collodom", Level: 1},
	} {
		if r.Description == "" {
			r.Description = "(room description not yet extracted from the original game)"
		}
		w.AddRoom(r)
	}
	w.Rooms[roomMisery].Visited = true
	return w
}

// Room IDs for CollodonsPile. Arbitrary/assigned by us (not extracted —
// the original's actual room-numbering scheme, if any, is unknown), just
// stable identifiers for this graph.
const (
	roomMisery RoomID = iota + 1
	roomSecundaPorta
	roomTrollwynd
	roomAgileStair
	roomMethos
	roomSothicComplex
	roomWolfdorp
	roomStings
	roomMorfang
	roomArrows
	roomNidus
	roomPilefoot
	roomPileCollodom
)
