package world

// LevelItem records that a named item (or obstacle/encounter marker) sits
// somewhere within a specific dungeon level, per the published maze-map
// poster included with Heavy on the Magick ("heavymap-levels1-2.jpg" in
// this repo — a Your Sinclair magazine map; see ../../CLAUDE.md for
// provenance). The map draws a hand-illustrated maze grid per level with
// item labels scattered across cells and thin lines marking open/closed
// passages between them.
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
// with a "+".
var LevelTwoItems = []LevelItem{
	{2, "Guards (2)"},
	{2, "Rock-Snake"},
	{2, "Lithic Key"},
	{2, "Alum Key"},
	{2, "One Magick Grail"},
	{2, "Bone"},
	{2, "Slug"},
	{2, "Egg"},
	{2, "Cabinet"},
	{2, "Toll"},
	{2, "Guards"},
	{2, "Copper Key"},
	{2, "Clasp (Tire)"},
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
