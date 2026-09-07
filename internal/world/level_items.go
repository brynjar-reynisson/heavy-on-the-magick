package world

// LevelItem records that a named item (or obstacle/encounter marker) sits
// somewhere within a specific dungeon level, per the published maze-map
// poster included with Heavy on the Magick ("heavymap-levels1-2.jpg" and
// its sibling "heavymap-levels3-4-poster.jpg" in this repo — a Your
// Sinclair magazine map; see ../../CLAUDE.md for provenance). The map
// draws a hand-illustrated maze grid per level with item labels scattered
// across cells and thin lines marking open/closed passages between them.
//
// Only the item-to-level association is captured here, not exact grid
// position or cell-by-cell wall connectivity: reading precise passage
// openings off a small, compressed, hand-drawn image is easy to get
// subtly wrong while still looking authoritative, so that data is
// deliberately NOT modeled as fact here rather than risk presenting a
// fabricated room graph as real. "Pile Collodom" appearing on this map
// (as "COLLODON / PILE", Level One) independently cross-confirms the
// room of the same name already sourced from the CASA walkthrough (see
// CollodonsPile) — good evidence both sources describe the same real
// game.
type LevelItem struct {
	Level int
	Name  string
}

// LevelOneItems are the item/obstacle labels legible on the map for Level
// One. Repeated labels are transcribed as read, not deduplicated — they
// may be distinct instances at different cells. CORRECTED UNDERSTANDING
// (round 31): "Toll" entries were originally assumed to be pickupable
// items; a later source (Spectrum Computing's instructions file - see
// ../../CLAUDE.md and internal/game's TollItem/payToll) confirms "Toll"
// actually marks a Toll DOOR location (needing a carried Bag of Gold to
// open), not an item — left in this list as-is (each entry still marks a
// real labeled cell on the map worth knowing about) rather than removed,
// but don't treat "Toll" here as something PICKUP would find.
var LevelOneItems = []LevelItem{
	{1, "Chroma Key"},
	{1, "Tar"},
	{1, "Bag"},
	{1, "Sword"},
	{1, "Bag"},
	{1, "Foot Garlic Bag"},
	{1, "Zinc Key"},
	{1, "Tin Key"},
	{1, "Bone"},
	{1, "Loaf"},
	{1, "Toll"},
	{1, "Shell (2 Jars)"},
	{1, "Toll"},
	{1, "Chest"},
	{1, "Pile Collodom"},
	{1, "Fire"},
	{1, "Cabinet"},
	{1, "Flask Cake"},
	{1, "Bond"},
}

// LevelTwoItems are the item/obstacle labels legible on the map for Level
// Two. "One Magick Grail" is annotated near the Alum Key/Guards cell on
// the source map — transcribed together since the map itself joins them
// with a "+". CORRECTED (round 149): re-reading the same source at a
// genuinely higher resolution than earlier rounds managed (see round
// 147/148's similar re-reads of this poster's OTHER half) shows two
// original transcription errors, now fixed: "Clasp (Tire)" is really
// "Clasp (Fire)" (matching the already-confirmed real Clasp/Fire
// mechanic — see world.Room.Fire's doc comment) and "One Magick Grail"
// is really "One Magick Grade" (matching Level Three's own "+ ONE
// MAGICK GRADE" cell below, and round 123's independent finding of the
// same "+ ONE MAGICK GRADE" marker).
var LevelTwoItems = []LevelItem{
	{2, "Guards (2)"},
	{2, "Rock-Snake"},
	{2, "Lithic Key"},
	{2, "Alum Key"},
	{2, "One Magick Grade"},
	{2, "Bone"},
	{2, "Slug"},
	{2, "Egg"},
	{2, "Cabinet"},
	{2, "Toll"},
	{2, "Guards"},
	{2, "Copper Key"},
	{2, "Clasp (Fire)"},
	{2, "Spell-Call"},
	{2, "Sign"},
	{2, "Grimoire Book"},
	{2, "Two Bags"},
	{2, "Nickel Key"},
	{2, "Bone"},
	{2, "Guards (1)"},
	{2, "Loaf"},
	{2, "Sun-Flower"},
	{2, "Toll"},
}

// LevelThreeItems and LevelFourItems are the item/obstacle labels
// legible on the OTHER half of the same official poster (round 149) -
// "heavymap-levels3-4-poster.jpg" in this repo, the sibling file to
// heavymap-levels1-2.jpg (source of LevelOneItems/LevelTwoItems
// above), read at full resolution for the first time this project has
// managed (matching round 147/148's similar re-read of the poster's
// other half). Same honest convention throughout: item-to-level
// association only, no grid position or connectivity, transcribed as
// read, not deduplicated.
//
// Real cross-confirmations found while transcribing: "Spell-
// Transfusion" (Level Three) matches the already-confirmed real
// TRANSFUSION spell and numbered_room_contents.go's own #29 "Scroll
// (TRANSFUSION spell)"; "Scroll Cauldron" (Level Three) independently
// corroborates #50's "Cauldron of cold iron (scroll inside)" from a
// completely different source (the numbered map poster); "Rabak"
// (Level Four) matches the already-placed "Doubt of Rabak" special
// room (level4_grid.go's D3); and Level Four's own "Skull"/"Head
// Bone"/"Iron Key" ... "Thigh Bone"/"Ulna Head" cluster independently
// corroborates round 139's CAULDRON, ACHAD ritual ingredients (Ulna,
// Thigh, Skull) all appearing together on the SAME level this
// project's own sourcing already pointed to. "+ One Magick Grade"
// (Level Three) is a SECOND real "+1 Grade" location distinct from
// Level Two's own (round 84's original find) - confirms this isn't a
// one-off marker, though neither is tied to a specific room yet.
var LevelThreeItems = []LevelItem{
	{3, "Limax Pellet"},
	{3, "Bronze Key"},
	{3, "Magnam Key"},
	{3, "Nougat"},
	{3, "Mantis"},
	{3, "Chest"},
	{3, "Rock"},
	{3, "Jar"},
	{3, "Chest"},
	{3, "Rock-Snake"},
	{3, "Clasp"},
	{3, "Rock"},
	{3, "Rock"},
	{3, "Cobalt Key"},
	{3, "Leaf Bag"},
	{3, "Brass Key"},
	{3, "Spell-Transfusion"},
	{3, "One Magick Grade"},
	{3, "Scroll Cauldron"},
	{3, "Water (Ball)"},
}

var LevelFourItems = []LevelItem{
	{4, "Bone"},
	{4, "Skull"},
	{4, "Head Bone"},
	{4, "Bone"},
	{4, "Bone"},
	{4, "Bone"},
	{4, "Toll"},
	{4, "Iron Key"},
	{4, "Rabak"},
	{4, "Rock"},
	{4, "Silver Nugget"},
	{4, "Thigh Bone"},
	{4, "Ulna Head"},
	{4, "Bone"},
	{4, "Chasm (Flask)"},
	{4, "Bag"},
	{4, "Medusa"},
}
