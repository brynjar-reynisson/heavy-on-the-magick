package world

// ZoneMonsterSighting records a monster type seen within a named zone on
// the computer-rendered grid map (see known_room_names.go's doc comment
// for the source, HeavyOnTheMagick_5.gif). Recorded at zone granularity,
// not per-cell: the source map places each icon in one specific lettered
// cell (e.g. a Werewolf at C3), but which cell that corresponds to in
// CollodonsPile's coarser room graph isn't known, so only "this zone has
// this monster type, this many times" is captured — deliberately not
// attached to any specific world.Room to avoid overwriting or
// contradicting monster data already sourced from the CASA walkthrough.
//
// Cyclops in the Nidus zone is a good independent cross-confirmation: it
// matches CollodonsPile's Nidus room (likely the same zone - this room
// was itself corrected from "Midus" to "Nidus" per the game's own
// vocabulary, see collodons_pile.go), which already has a Cyclops from
// the walkthrough source.
//
// Monster names here use this project's own canonical spelling, same as
// every other Room.Monster value - the source map's own legend actually
// glossed this creature "wraith", corrected project-wide to "Vampire" in
// round 74 (see Level1Grid's doc comment for the full reasoning).
//
// ROUND 177 CORRECTION: Methos's entry below was reverted from
// "Vampire" back to "Wraith" - a full frame-by-frame review of real
// gameplay footage shows Methos's own live combat text literally reads
// "WRAITH IS DEAD", while a SEPARATE encounter (Morfang) shows real
// combat text "VAMPIRE ATTACKS!" for a genuinely different monster -
// direct evidence "Wraith" and "Vampire" are two real, distinct
// creatures, not one creature under two source-dependent names as
// round 74 assumed. This reopens whether round 74's OTHER 6 renamed
// placements (Level1Grid F2/G1/G2/H1, Level2Grid A5, Level4Grid A6)
// were also over-corrected - none of those have been re-checked
// against real gameplay text yet, so they're deliberately left as
// "Vampire" rather than reverted on inference alone. The "Wraithvale"
// zone's own "Vampire" sighting below is worth particular suspicion
// (a zone literally named after "Wraith" reporting a "Vampire"
// sighting) but is likewise left unchanged pending real evidence - see
// CLAUDE.md's "Open next steps".
type ZoneMonsterSighting struct {
	Zone    string
	Monster string
	Count   int
}

var ZoneMonsterSightings = []ZoneMonsterSighting{
	{Zone: "Wolfdorp", Monster: "Ghost", Count: 2},
	{Zone: "Wolfdorp", Monster: "Werewolf", Count: 2},
	{Zone: "Morfang", Monster: "Vampire", Count: 3},
	{Zone: "Nidus", Monster: "Cyclops", Count: 1},
	{Zone: "Wraithvale", Monster: "Vampire", Count: 1},
	{Zone: "Room of Icthys", Monster: "Slug", Count: 1},
	{Zone: "Sothic Complex", Monster: "Ghost", Count: 1},
	{Zone: "Gorburg", Monster: "Wyvern", Count: 1},
	{Zone: "Gorburg", Monster: "Ghost", Count: 2},
	{Zone: "Trollwynd", Monster: "Troll", Count: 4},
	{Zone: "Rook of Hydra", Monster: "Wyvern", Count: 1},
	{Zone: "Wormring", Monster: "Wyvern", Count: 4},
	{Zone: "Doubt of Rabak", Monster: "Vampire", Count: 1},
	{Zone: "Methos", Monster: "Wraith", Count: 1},
	{Zone: "The Pit", Monster: "Medusa", Count: 1},
}
