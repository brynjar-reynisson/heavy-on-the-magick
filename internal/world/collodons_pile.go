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
// Trollwynd's monster was left generically named "monster" for many
// rounds even after zone_monsters.go's independent ZoneMonsterSightings
// ("Trollwynd: Troll x4") was cross-checked against Level3Grid's own
// data and found to match EXACTLY (4 tight-crop-verified Trolls at
// C4/C6/E7/F8, all within the Trollwynd zone — see the Methos/Vampire
// note below, which used this same match as supporting evidence but
// never went back and applied it here). Round 101: named it "Troll"
// accordingly — not a new source, just finally acting on already-
// verified evidence sitting a few paragraphs below this one.
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
// locations") rather than one specific room, so they were deliberately
// NOT placed then — assigning them to a guessed room would have been
// fabrication.
//
// Methos's Vampire (round 71, renamed from "Wraith" in round 74 - see
// Level1Grid's doc comment): zone_monsters.go's ZoneMonsterSightings
// (real data, sourced back in round 12 but never fully cross-checked
// against every already-real room until now) records "Methos: Wraith
// x1" (using this project's since-corrected name for the creature) -
// and Methos was the one CollodonsPile room in that list still
// missing a monster entirely. Cross-checking the OTHER sightings
// against already-shipped per-cell data found strong, exact
// corroboration everywhere else (Trollwynd's "Troll x4" matches
// Level3Grid's 4 tight-crop-verified Trolls exactly; Gorburg's
// "Wyvern x1, Ghost x2" matches Level3Grid's Gorburg-zone placements
// exactly; Wolfdorp's "Ghost x2, Werewolf x2" matches Level1Grid's
// placements exactly) - real, strong validation that both this
// independent source and the per-cell tight-crop work agree. Given
// that track record, Methos's still-missing monster was trustworthy
// too, and unlike the isolated-cell finds elsewhere in this project,
// Methos is already a real, connected, playable CollodonsPile room -
// this makes it a genuine, reachable combat encounter, not just
// recorded data.
//
// Morfang's Vampire (round 106): zone_monsters.go records "Morfang:
// Vampire x3" — at the time of Methos's fix above (round 72's
// writeup), this was deliberately left untouched, flagged as "a
// near-miss" because Level1Grid separately ships 4 Vampires (F2/G1/G2/
// H1), not 3, an exact-count mismatch unlike Trollwynd/Gorburg/
// Wolfdorp's clean matches. Revisited: that earlier caution conflated
// two different questions. The COUNT genuinely doesn't reconcile (real,
// still-open, not resolved here) — but CollodonsPile's Room.Monster
// field only records WHICH species is present, not how many, and both
// independent sources agree unambiguously on that: Vampire. The count
// mismatch has no bearing on that simpler question, so leaving Morfang
// with no monster at all (the status quo) was needlessly conservative,
// not actually required by the discrepancy. Same MonsterHealth=2
// placeholder as Level1Grid's own Vampire cells (not a new number).
//
// HasChest placements (round 78): a fresh, more targeted CASA re-read
// (asking specifically whether "EXAMINE X, pick up Y" is a consistent
// pattern across the whole walkthrough, not just Wolfdorp) confirmed
// "EXAMINE CHEST" is a real, distinct container phrase from "EXAMINE
// TABLE" - used before "Pick up GARLIC" in Wolfdorp and "Pick up SLAT"
// in Morfang specifically (every other EXAMINE-before-pickup instance
// in the walkthrough uses TABLE, already modeled). Added
// world.Room.HasChest, wired into examine() and describeCurrentRoom()
// the same way HasTable already is - same honest "only these 2
// specifically confirmed rooms" convention, not assumed elsewhere.
//
// CORRECTED (round 82): round 81 added "Key" to Wolfdorp's Items,
// citing "(Wolfdorp on level 1) ... EXAMINE TABLE, Pick up KEY" - but
// that phrasing came from an AI-summarized LIST ("every EXAMINE-before-
// pickup instance"), not a direct quote, and it was wrong. A follow-up
// round asked for the raw literal source text character-for-character
// around this exact passage and got: "...(Wolfdorp on level 1),
// EXAMINE TABLE, Pick up LOAF, W, \"DOOR LUNACY\" (door opens), N, DROP
// CLASP, Pick up KEY, SW, W, SW, S, S, NW (Room of stings...". Read
// correctly, "Pick up KEY" happens 2 moves and a door-password AWAY
// from Wolfdorp, in an unnamed intermediate room - not at Wolfdorp
// itself. This is exactly the same "unnamed room between two named
// ones" trap round 64 already caught once (the original Pilefoot/DROP
// KEY case) - round 81 fell into a variant of it by trusting a
// summarized list instead of checking the raw quote. "Key" has been
// REMOVED from Wolfdorp's Items; Room of Stings' TollItem "Key" is
// once again honestly unplaced/unsourced within CollodonsPile, the
// same status it held from round 64 through round 80. Lesson for
// future rounds: when a categorized/summarized answer places a fact in
// a named room, verify against the RAW literal source text before
// shipping - a summary can silently misattribute an action to whichever
// room name happens to appear nearest it in the text, even across
// several intervening moves.
//
// TollItem placements (round 64): the same fresh, more detailed CASA
// walkthrough re-read found the real drop-to-open-door mechanic (see
// world.Room.TollItem's doc comment) recurring at 3 more rooms, each
// individually confirmed under that room's own paragraph heading (not
// just nearby text): Room of Stings needs a Key dropped, Morfang needs
// the Bag (already placed in Wolfdorp — a real, satisfying pickup-then-
// use chain), and Room of Arrows needs the Slat (already placed in
// Morfang — likewise). A fourth apparent "DROP KEY" instance near
// Pilefoot in the raw text turned out, on closer checking, to be listed
// under a different, unidentified room's heading, not Pilefoot's own -
// left unplaced rather than guess which room it really belongs to.
//
// CORRECTED (round 63): a fresh, more detailed re-read of the same CASA
// walkthrough resolved 2 of those 3 - it gives Nougat and a Scroll both
// specifically in the Trollwynd area, and a SECOND separate Scroll
// specifically in Sothic Complex (i.e. 2 real Scroll instances at 2
// named rooms, matching the original pass's "two instances" hedge -
// that wasn't vague after all, just under-extracted the first time).
// Key remains genuinely spread across 4 different rooms with no single
// location, so it's still deliberately unplaced. The same re-read also
// surfaced a real, previously-unknown mechanic: Werewolves are
// "killable by walking through after dropping NOUGAT" - see
// game.checkNougatWerewolf, wired in for the first time this round.
//
// Trollwynd's Clasp/Scroll cross-referenced (round 87): the numbered
// map poster's key list has #21 "Cabinet (clasp - Salamander charm)"
// and #22 "Scroll (CALL spell)", both tight-crop-confirmed to sit
// within the Trollwynd zone banner, right alongside #24 (Nougat,
// already placed here). Real evidence this Clasp and this Scroll ARE
// those exact numbered-map items, not just coincidentally-named
// duplicates - see game.Handle's CALL doc comment for how this
// sharpens (without fully resolving) the CALL spell's honest stub.
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
//
// Wolfdorp's "Sword" (round 52): cross-referenced from TWO independent
// sources. The numbered map poster's key list gives room #65 as "Rock,
// two stalagmites, stalactite, sword", and #65 sits within the
// "WOLFDORP" banner-labeled cluster on that same poster's Level 1 grid
// (tight-cropped and visually confirmed, not guessed). Separately,
// level_items.go's LevelOneItems (from the OTHER poster,
// heavymap-levels1-2.jpg, extracted in an earlier round) independently
// lists a "Sword" on Level 1 with no room precision. Two unrelated
// sources agreeing Level 1 has a sword, one of them at zone-level
// confidence for Wolfdorp specifically, matches the Mantis/Belezbar
// precedent (Level3Grid) closely enough to place it here. "Sword" is
// Astarot's confirmed real Charm (magic.Demons) — this makes Astarot's
// invocation reachable in real gameplay for the first time.
//
// Room of Misery's "Poison-smeared book" (round 53): the numbered map
// poster labels this exact room "START / Room of Misery / 1, 2" - i.e.
// Room of Misery IS numbered cells #1 and #2 on that poster, no
// zone-level guessing needed at all (the strongest-confidence item
// placement in this file). #1 is "Grimoire" - already placed here
// independently via the CASA walkthrough, an exact cross-source match
// that validates this numbered-cell reading. #2 is "Poison-smeared
// book", not previously placed anywhere; added here on the same
// footing as the already-confirmed Grimoire.
//
// Sothic Complex's "Sunflower" (round 77): same numbered-map-plus-zone-
// banner cross-reference method as Wolfdorp's Sword above. The poster's
// key list gives room #7 as "Chest (Sunflower)", and #7 sits directly
// within the "SOTHIC COMPLEX" banner-labeled cell cluster on that same
// poster's Level 2 grid (tight-cropped and visually confirmed - #7 is
// NOT inside the neighboring "KITCHEN OF AI" banner's cluster, checked
// directly since the two banners sit close together). "Sunflower" is
// Magot's confirmed real Charm (magic.Demons) - this makes Magot's new
// "MAGOT, <object>" locate ability (see game.magotLocate) reachable in
// real gameplay for the first time, the same way Sword's placement here
// made Astarot's teleport reachable.
func CollodonsPile() *World {
	w := New(roomMisery)
	for _, r := range []*Room{
		{ID: roomMisery, Name: "Room of Misery", Level: 2, Exits: map[Direction]RoomID{East: roomSecundaPorta}, Items: []string{"Grimoire", "Poison-smeared book"}, HasTable: true},
		{ID: roomSecundaPorta, Name: "Secunda Porta", Level: 2, Exits: map[Direction]RoomID{North: roomTrollwynd}, DoorPasswords: []string{"SILENCE"}},
		{ID: roomTrollwynd, Name: "Trollwynd", Level: 3, Exits: map[Direction]RoomID{North: roomAgileStair, South: roomSothicComplex}, Monster: "Troll", MonsterHealth: 3, Items: []string{"Clasp", "Nougat", "Scroll"}, HasTable: true},
		{ID: roomAgileStair, Name: "Agile Stair", Level: 4, Exits: map[Direction]RoomID{SouthEast: roomMethos}},
		{ID: roomMethos, Name: "Methos", Level: 4, Exits: map[Direction]RoomID{South: roomSothicComplex}, Items: []string{"Nugget"}, HasTable: true, Monster: "Vampire", MonsterHealth: 2},
		{ID: roomSothicComplex, Name: "Sothic Complex", Level: 2, Exits: map[Direction]RoomID{South: roomWolfdorp}, Items: []string{"Scroll", "Sunflower"}, HasTable: true},
		{ID: roomWolfdorp, Name: "Wolfdorp", Level: 1, Exits: map[Direction]RoomID{NorthWest: roomStings}, DoorPasswords: []string{"WOLF", "LUNACY"}, Items: []string{"Garlic", "Bag", "Loaf", "Sword"}, HasTable: true, HasChest: true},
		{ID: roomStings, Name: "Room of Stings", Level: 1, Exits: map[Direction]RoomID{North: roomMorfang}, TollItem: "Key", HasTable: true},
		{ID: roomMorfang, Name: "Morfang", Level: 1, Exits: map[Direction]RoomID{East: roomArrows}, Monster: "Vampire", MonsterHealth: 2, Items: []string{"Slat"}, TollItem: "Bag", HasTable: true, HasChest: true},
		{ID: roomArrows, Name: "Room of Arrows", Level: 1, Exits: map[Direction]RoomID{East: roomNidus, North: roomWolfdorp}, TollItem: "Slat", HasTable: true},
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
